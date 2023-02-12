package integration

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"

	"github.com/morozovcookie/ogen/conv"
	ht "github.com/morozovcookie/ogen/http"
	api "github.com/morozovcookie/ogen/internal/integration/sample_api"
	"github.com/morozovcookie/ogen/internal/integration/techempower"
	"github.com/morozovcookie/ogen/json"
)

func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	return l
}

type RPS struct {
	count int64
	start time.Time
}

func (r *RPS) Inc() {
	atomic.AddInt64(&r.count, 1)
}

func (r *RPS) Report(b *testing.B) {
	sec := time.Since(r.start).Seconds()
	perSec := float64(atomic.LoadInt64(&r.count)) / sec
	b.ReportMetric(perSec, "req/s")
}

func newRPS() *RPS {
	return &RPS{
		start: time.Now(),
	}
}

func BenchmarkValidation(b *testing.B) {
	pet := &api.Pet{
		Kind: api.PetKindSmol,
		Name: "Foo Bar",
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := pet.Validate(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchClient() *http.Client {
	const (
		maxConnsPerHost     = 20
		maxIdleConns        = maxConnsPerHost
		maxIdleConnsPerHost = maxIdleConns
	)
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          maxIdleConns,
			MaxIdleConnsPerHost:   maxIdleConnsPerHost,
			MaxConnsPerHost:       maxConnsPerHost,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func BenchmarkIntegration(b *testing.B) {
	b.Run("Baseline", func(b *testing.B) {
		// Use baseline implementation to measure framework overhead.
		b.Run("Std", func(b *testing.B) {
			data := []byte(`Hello, world!`)
			b.SetBytes(int64(len(data)))
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(data)
			}))
			defer s.Close()

			client := s.Client()

			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if err := func() error {
						res, err := client.Get(s.URL)
						if err != nil {
							return err
						}
						defer func() {
							_ = res.Body.Close()
						}()
						if _, err := io.ReadAll(res.Body); err != nil {
							return err
						}
						if res.StatusCode != http.StatusOK {
							return errors.Errorf("code: %d", res.StatusCode)
						}

						return nil
					}(); err != nil {
						b.Error(err)
					}
				}
			})
		})
		b.Run("Fasthttp", func(b *testing.B) {
			done := make(chan struct{})
			defer func() { <-done }()

			ln := newLocalListener()
			defer func() { _ = ln.Close() }()

			go func() {
				defer close(done)
				if err := fasthttp.Serve(ln, func(ctx *fasthttp.RequestCtx) {
					_, _ = ctx.WriteString("Hello, world!")
				}); err != nil {
					b.Error(err)
				}
			}()

			c := &fasthttp.Client{}
			u := (&url.URL{
				Host:   ln.Addr().String(),
				Scheme: "http",
			}).String()

			b.ResetTimer()
			b.ReportAllocs()
			b.SetBytes(int64(len("Hello, world!")))
			rps := newRPS()
			defer rps.Report(b)
			b.RunParallel(func(pb *testing.PB) {
				var dst []byte
				for pb.Next() {
					code, result, err := c.Get(dst, u)
					if err != nil {
						b.Error(err)
						return
					}
					if code != http.StatusOK {
						b.Errorf("bad code %d:", code)
						return
					}

					// Reusing buffer.
					dst = result[:0]
					rps.Inc()
				}
			})
		})
	})

	b.Run("Manual", func(b *testing.B) {
		// Test with some manual optimizations.
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			e := jx.GetEncoder()
			e.ObjStart()
			e.FieldStart("message")
			e.Str("Hello, world!")
			e.ObjEnd()
			if _, err := e.WriteTo(w); err != nil {
				b.Error(err)
			}
			jx.PutEncoder(e)
		}))
		defer s.Close()

		ctx := context.Background()
		client := benchClient()
		b.Run("JSON", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			u, err := url.Parse(s.URL)
			require.NoError(b, err)
			rps := newRPS()
			defer rps.Report(b)

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					req, err := ht.NewRequest(ctx, http.MethodGet, u, nil)
					if err != nil {
						b.Error(err)
						break
					}
					res, err := client.Do(req)
					if err != nil {
						b.Error(err)
						break
					}
					_, _ = io.Copy(io.Discard, res.Body)
					_ = res.Body.Close()
					if res.StatusCode != http.StatusOK {
						b.Error(res.StatusCode)
						break
					}
					rps.Inc()
				}
			})
		})
	})

	b.Run("TechEmpower", func(b *testing.B) {
		// Using TechEmpower as most popular general purpose framework benchmark.
		// https://github.com/TechEmpower/FrameworkBenchmarks/wiki/Project-Information-Framework-Tests-Overview#test-types

		srv := techEmpowerServer{}
		h, err := techempower.NewServer(srv)
		require.NoError(b, err)

		s := httptest.NewServer(h)
		defer s.Close()

		httpClient := benchClient()
		client, err := techempower.NewClient(s.URL,
			techempower.WithClient(httpClient),
			techempower.WithTracerProvider(trace.NewNoopTracerProvider()),
		)
		require.NoError(b, err)
		ctx := context.Background()

		b.Run("JSON", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			rps := newRPS()
			defer rps.Report(b)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					hw, err := client.JSON(ctx)
					if err != nil {
						b.Fatal(err)
						return
					}
					if hw.Message != "Hello, world!" {
						b.Error("mismatch")
					}
					rps.Inc()
				}
			})
		})
		b.Run("OnlyHandler", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					hw, err := srv.JSON(ctx)
					if err != nil {
						b.Error(err)
						return
					}
					if hw.Message != "Hello, world!" {
						b.Error("mismatch")
					}
				}
			})
		})
	})
}

func BenchmarkJSON(b *testing.B) {
	b.Run("TechEmpower", func(b *testing.B) {
		b.Run("HelloWorld", func(b *testing.B) {
			h := &techempower.HelloWorld{
				Message: "Hello, world!",
			}
			data := json.Encode(h)
			dataBytes := int64(len(data))

			b.Run("Encode", func(b *testing.B) {
				var e jx.Encoder
				b.ReportAllocs()
				b.SetBytes(dataBytes)

				for i := 0; i < b.N; i++ {
					e.Reset()
					h.Encode(&e)
				}
			})
			b.Run("Decode", func(b *testing.B) {
				var v techempower.HelloWorld
				b.ReportAllocs()
				b.SetBytes(dataBytes)
				d := jx.GetDecoder()

				for i := 0; i < b.N; i++ {
					d.ResetBytes(data)
					if err := v.Decode(d); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
		b.Run("WorldObject", func(b *testing.B) {
			h := &techempower.WorldObject{
				ID:           367297,
				RandomNumber: 4761696123,
			}
			data := json.Encode(h)
			dataBytes := int64(len(data))

			b.Run("Encode", func(b *testing.B) {
				var e jx.Encoder
				b.ReportAllocs()
				b.SetBytes(dataBytes)

				for i := 0; i < b.N; i++ {
					e.Reset()
					h.Encode(&e)
				}
			})
			b.Run("Decode", func(b *testing.B) {
				var v techempower.WorldObject
				b.ReportAllocs()
				b.SetBytes(dataBytes)
				j := jx.GetDecoder()

				for i := 0; i < b.N; i++ {
					j.ResetBytes(data)
					if err := v.Decode(j); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	})
	b.Run("Sample", func(b *testing.B) {
		b.Run("Pet", func(b *testing.B) {
			date := time.Date(2011, 10, 10, 7, 12, 34, 4125, time.UTC)
			pet := &api.Pet{
				Primary:  nil,
				ID:       42,
				UniqueID: uuid.New(),
				Name:     "SomePet",
				Type:     api.NewOptPetType(api.PetTypeFofa),
				Kind:     api.PetKindSmol,
				Tag:      api.NewOptUUID(uuid.New()),
				IP:       netip.AddrFrom4([4]byte{1, 1, 1, 1}),
				IPV4:     netip.AddrFrom4([4]byte{8, 8, 8, 8}),
				IPV6:     netip.AddrFrom16([16]byte{15: 1}),
				URI: func() url.URL {
					u, _ := url.Parse("https://google.com")
					return *u
				}(),
				Birthday:     conv.Date(date),
				Rate:         time.Second,
				Nickname:     api.NewNilString("Nick"),
				NullStr:      api.NewOptNilString("Bar"),
				Friends:      nil,
				Next:         api.OptData{},
				TestInteger1: api.NewOptInt(10),
				TestFloat1:   api.NewOptFloat64(1.0),
				TestArray1:   nil,
				TestDate:     api.NewOptDate(conv.Date(date)),
				TestDuration: api.NewOptDuration(time.Minute),
				TestTime:     api.NewOptTime(conv.Time(date)),
				TestDateTime: api.NewOptDateTime(conv.DateTime(date)),
			}
			data := json.Encode(pet)
			dataBytes := int64(len(data))
			b.Run("Encode", func(b *testing.B) {
				var e jx.Encoder

				b.ReportAllocs()
				b.SetBytes(dataBytes)
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					e.Reset()
					pet.Encode(&e)
				}
			})
			b.Run("Decode", func(b *testing.B) {
				var d jx.Decoder

				b.ReportAllocs()
				b.SetBytes(dataBytes)
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					d.ResetBytes(data)
					if err := pet.Decode(&d); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
		b.Run("PetType", func(b *testing.B) {
			pet := api.PetTypeFofa
			data := json.Encode(pet)
			dataBytes := int64(len(data))
			b.Run("Encode", func(b *testing.B) {
				var e jx.Encoder

				b.ReportAllocs()
				b.SetBytes(dataBytes)
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					e.Reset()
					pet.Encode(&e)
				}
			})
			b.Run("Decode", func(b *testing.B) {
				var d jx.Decoder

				b.ReportAllocs()
				b.SetBytes(dataBytes)
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					d.ResetBytes(data)
					if err := pet.Decode(&d); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	})
}
