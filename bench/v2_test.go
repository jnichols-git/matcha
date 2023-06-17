package bench

import (
	"math/rand"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/cloudretic/matcha/pkg/route"
	"github.com/cloudretic/matcha/pkg/router"
)

/*
 * ===== MockBoards V2 API =====
 *
 * Use Case: MockBoards is a (fake) forum website setting up their API using Matcha. They're updating their API,
 * but to maintain backwards compatibility, they're attaching v2 to a new endpoint. For the purposes of this benchmark,
 * both APIs are the same--this is just a test of the performance of subrouters.
 *
 * ===== Using This Benchmark =====
 *
 * Run the API benchmark, and the offset benchmark from v1_test.go.
 * For sequential results, subtract the offset from the X/op values.
 * For concurrent results, divide the X/op values by 10, then subtract the offset.
 */

func choosev2() benchRoute {
	idx := rand.Int() % len(apiRoutes)
	choice := apiRoutes[idx]
	v2 := rand.Int() % 2
	if v2 == 0 {
		choice.testPath = "/v2" + choice.testPath
	}
	return choice
}

// Just to check!
func TestAPIv2(t *testing.T) {
	v1 := router.Default()
	v2 := router.Default()
	for _, tr := range apiRoutes {
		r := route.Declare(tr.method, tr.path)
		for _, mw := range tr.mws {
			r.Attach(mw)
		}
		for _, rq := range tr.rqs {
			r.Require(rq)
		}
		v1.HandleRouteFunc(r, handleOK)
		v2.HandleRouteFunc(r, handleOK)
	}
	v1.Mount("/v2", v2)
	w := httptest.NewRecorder()
	for i := 0; i < len(apiRoutes); i++ {
		br := apiRoutes[i]
		req := httptest.NewRequest(br.method, br.testPath, nil)
		req.Header.Set("X-Platform-User-ID", "jnichols")
		v1.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatal(br.method, br.path, br.testPath, w.Code)
		}
		req = httptest.NewRequest(br.method, "/v2"+br.testPath, nil)
		req.Header.Set("X-Platform-User-ID", "jnichols")
		v1.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatal(br.method, br.path, br.testPath, w.Code)
		}
	}
}

func BenchmarkAPIv2(b *testing.B) {
	v1 := router.Default()
	v2 := router.Default()
	for _, tr := range apiRoutes {
		r, err := route.New(tr.method, tr.path)
		if err != nil {
			b.Fatal(err)
		}
		for _, mw := range tr.mws {
			r.Attach(mw)
		}
		for _, rq := range tr.rqs {
			r.Require(rq)
		}
		v1.HandleRouteFunc(r, handleOK)
		v2.HandleRouteFunc(r, handleOK)
	}
	v1.Mount("/v2", v2)
	b.Run(b.Name()+"-sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w := httptest.NewRecorder()
			br := choosev2()
			req := httptest.NewRequest(br.method, br.testPath, nil)
			req.Header.Set("X-Platform-User-ID", "jnichols")
			v1.ServeHTTP(w, req)
		}
	})
	b.Run(b.Name()+"-concurrent-10", func(b *testing.B) {
		wg := &sync.WaitGroup{}
		for i := 0; i < b.N; i++ {
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					w := httptest.NewRecorder()
					br := choosev2()
					req := httptest.NewRequest(br.method, br.testPath, nil)
					req.Header.Set("X-Platform-User-ID", "jnichols")
					v1.ServeHTTP(w, req)
					wg.Done()
				}()
			}
			wg.Wait()
		}
	})
}

// BenchmarkOffset is used to calculate the performance offset of generating test values
// for benchmarks.
// This is important; creating new requests is expensive
func BenchmarkOffsetv2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = httptest.NewRecorder()
		br := choosev2()
		req := httptest.NewRequest(br.method, br.testPath, nil)
		req.Header.Set("X-Platform-User-ID", "jnichols")
	}
}
