package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func main() {
	// Создаем команду, например, "nix build"
	cmd := exec.Command("/nix/store/q3d20mdmc7wpl4jam09nqfr7dqss6ar1-nix-2.27.0pre19700101_dirty/bin/nix", "build", "--no-link", "--impure", "nixpkgs#hello")
	cmd.Env = []string{"TERM=xterm-256color"}
	// Создаем псевдотерминал
	pty, err := pty.Start(cmd)
	if err != nil {
		log.Fatalf("Ошибка при создании псевдотерминала: %v", err)
	}
	defer pty.Close()

	// Получаем файловый дескриптор псевдотерминала
	fd := pty.Fd()

	// Выводим файловый дескриптор для информации
	fmt.Printf("Файловый дескриптор псевдотерминала: %d\n", fd)

	// Копируем вывод псевдотерминала в стандартный вывод
	go func() {
		_, err := io.Copy(os.Stdout, pty)
		if err != nil {
			log.Printf("Ошибка при выводе на stdout: %v\n", err)
		}
		fmt.Println("Done")
	}()

	// Ждем завершения команды
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Команда завершилась с ошибкой: %v", err)
	}
	fmt.Println("\nКоманда завершена успешно")
	time.Sleep(10 * time.Second)
}
