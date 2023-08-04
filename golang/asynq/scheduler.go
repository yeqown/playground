package main

import (
	"context"

	"fmt"
	"github.com/hibiken/asynq"
)

func main() {
	scheduler := asynq.NewScheduler(getRedisConnOpt(), &asynq.SchedulerOpts{
		Logger:   nil,
		LogLevel: 0,
		Location: nil,
		PreEnqueueFunc: func(task *asynq.Task, opts []asynq.Option) {
			fmt.Println("PreEnqueueFunc")
			fmt.Printf("task playload: %s\n", task.Payload())

			// ResultWriter is nil
			// fmt.Printf("task id: %+v\n", task.ResultWriter().TaskID())
			// change task instance
		},
		PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
			fmt.Println("PostEnqueueFunc")
			fmt.Printf("info: %+v\n", info)
		},
	})
	_, err := scheduler.Register("@every 5s", asynq.NewTask("send_email", []byte("xxx@gmail.com")))
	if err != nil {
		panic(err)
	}

	server := asynq.NewServer(getRedisConnOpt(), asynq.Config{})
	serverMux := asynq.NewServeMux()
	serverMux.HandleFunc("send_email", func(ctx context.Context, task *asynq.Task) error {
		fmt.Printf("send_email called: %s\n", task.ResultWriter().TaskID())
		return nil
	})

	go func() {
		if err := server.Run(serverMux); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := scheduler.Run(); err != nil {
			panic(err)
		}
	}()

	// block forever
	select {}
}

func getRedisConnOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}
