#!/bin/bash

echo "=== Тестирование mysh ==="

# Функция для тестирования вывода команд
test_command() {
    local test_name="$1"
    local command="$2"
    local expected_output="$3"

    echo ""
    echo "=== $test_name ==="
    echo "Команда: $command"

    # Запускаем mysh с командой и захватываем вывод
    local actual_output
    actual_output=$({ sleep 0.1; echo "$command"; sleep 0.1; } | ./mysh)

    echo "--- Ожидаемый вывод: ---"
    echo "$expected_output"
    echo "--- Фактический вывод: ---"
    echo "$actual_output"

    if [ "$actual_output" = "$expected_output" ]; then
        echo "✓ Вывод соответствует"
    else
        echo "✗ Вывод не соответствует"
    fi
}

# Функция для тестирования редиректов в файлы
test_redirect() {
    local test_name="$1"
    local command="$2"
    local test_file="$3"
    local expected_content="$4"

    echo ""
    echo "=== $test_name ==="
    echo "Команда: $command"

    # Очищаем файл
#    rm -f "$test_file"

    # Запускаем mysh с командой редиректа
    { sleep 0.1; echo "$command"; sleep 0.1; } | ./mysh

    echo "--- Ожидаемое содержимое: ---"
    echo "$expected_content"
    echo "--- Фактическое содержимое: ---"
    if [ -f "$test_file" ]; then
        cat "$test_file"
        local actual_content
        actual_content=$(cat "$test_file")
        if [ "$actual_content" = "$expected_content" ]; then
            echo "✓ Файл создан и содержимое соответствует"
        else
            echo "✗ Содержимое не соответствует"
        fi
    else
        echo "Файл не создан"
        echo "✗ Файл не создан"
    fi
}

# Очистка
cleanup() {
    rm -f test1.txt test2.txt test3.txt test4.txt test5.txt
}

trap cleanup EXIT

echo "=== Базовые команды ==="
test_command "Тест 1: pwd" "pwd" "/home/kromka/WB_L2/L2_15"
test_command "Тест 2: echo" "echo hello world" "hello world"

echo ""
echo "=== Пайплайны ==="
test_command "Тест 3: Простой пайплайн" "echo hello world | grep hello" "hello world"
test_command "Тест 4: Пайплайн с wc" "echo line1 | wc -c" "6"

echo ""
echo "=== Логические операторы ==="
test_command "Тест 5: AND успех" "pwd && echo success" "/home/kromka/WB_L2/L2_15
success"
test_command "Тест 6: OR с ошибкой" "cd nonexistent || echo fallback" "fallback"

echo ""
echo "=== Редиректы в файлы ==="
test_redirect "Тест 1: Редирект вывода" "pwd > test1.txt" "test1.txt" "/home/kromka/WB_L2/L2_15"
test_redirect "Тест 2: Дополнение файла" "echo line1 >> test2.txt" "test2.txt" "line1"
test_redirect "Тест 3: Дополнение файла 2" "echo line2 >> test2.txt" "test2.txt" "line1
line2"

echo ""
echo "=== Тестирование завершено ==="