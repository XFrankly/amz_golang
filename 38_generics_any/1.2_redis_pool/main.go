package main

// func makeConnPool() {
// 	const minIdleConns = 10

// 	var (
// 		wg         sync.WaitGroup
// 		closedChan = make(chan struct{})
// 	)
// 	wg.Add(minIdleConns)
// 	connPool := pool.NewConnPool(&redis.Options{
// 		Dialer: func(ctx context.Context) (net.Conn, error) {
// 			wg.Done()
// 			<-closedChan
// 			return &net.TCPConn{}, nil
// 		},
// 		PoolSize:        10,
// 		PoolTimeout:     time.Hour,
// 		ConnMaxIdleTime: time.Millisecond,
// 		MinIdleConns:    minIdleConns,
// 	})
// 	wg.Wait()
// 	close(closedChan)
// }
// func main() {
// 	makeConnPool()
// }
