#!/bin/bash

# Тестовый скрипт для сравнения my_grep с GNU grep

# Создаем тестовый файл
cat > test_input.txt << EOF
apple line 1
banana line 2
cherry line 3
date line 4
elderberry line 5
fig line 6
grape line 7
honeydew line 8
apple line 9
banana line 10
APPLE line 11
Banana line 12
line with multiple apple words
just another line
final line with APPLE
EOF

echo "=== Тестовый файл создан ==="
cat test_input.txt

echo -e "\n=== Тест 1: Простой поиск ==="
echo "GNU grep 'apple':"
grep 'apple' test_input.txt > grep_output1.txt
cat grep_output1.txt

echo -e "\nMy grep 'apple':"
./my_grep 'apple' test_input.txt > my_grep_output1.txt
cat my_grep_output1.txt

echo -e "\nСравнение:"
if diff grep_output1.txt my_grep_output1.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 2: Игнорирование регистра (-i) ==="
echo "GNU grep -i 'apple':"
grep -i 'apple' test_input.txt > grep_output2.txt
cat grep_output2.txt

echo -e "\nMy grep -i 'apple':"
./my_grep -i 'apple' test_input.txt > my_grep_output2.txt
cat my_grep_output2.txt

echo -e "\nСравнение:"
if diff grep_output2.txt my_grep_output2.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 3: Инвертированный поиск (-v) ==="
echo "GNU grep -v 'apple':"
grep -v 'apple' test_input.txt > grep_output3.txt
cat grep_output3.txt | head -10 # выводим первые 10 строк

echo -e "\nMy grep -v 'apple':"
./my_grep -v 'apple' test_input.txt > my_grep_output3.txt
cat my_grep_output3.txt | head -10

echo -e "\nСравнение:"
if diff grep_output3.txt my_grep_output3.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 4: Номера строк (-n) ==="
echo "GNU grep -n 'apple':"
grep -n 'apple' test_input.txt > grep_output4.txt
cat grep_output4.txt

echo -e "\nMy grep -n 'apple':"
./my_grep -n 'apple' test_input.txt > my_grep_output4.txt
cat my_grep_output4.txt

echo -e "\nСравнение:"
if diff grep_output4.txt my_grep_output4.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 5: Только количество (-c) ==="
echo "GNU grep -c 'apple':"
grep -c 'apple' test_input.txt > grep_output5.txt
cat grep_output5.txt

echo -e "\nMy grep -c 'apple':"
./my_grep -c 'apple' test_input.txt > my_grep_output5.txt
cat my_grep_output5.txt

echo -e "\nСравнение:"
if diff grep_output5.txt my_grep_output5.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 6: Фиксированная строка (-F) ==="
echo "GNU grep -F 'apple':"
grep -F 'apple' test_input.txt > grep_output6.txt
cat grep_output6.txt

echo -e "\nMy grep -F 'apple':"
./my_grep -F 'apple' test_input.txt > my_grep_output6.txt
cat my_grep_output6.txt

echo -e "\nСравнение:"
if diff grep_output6.txt my_grep_output6.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 7: Контекст после (-A 2) ==="
echo "GNU grep -A 2 'cherry':"
grep -A 2 'cherry' test_input.txt > grep_output7.txt
cat grep_output7.txt

echo -e "\nMy grep -A 2 'cherry':"
./my_grep -A 2 'cherry' test_input.txt > my_grep_output7.txt
cat my_grep_output7.txt

echo -e "\nСравнение:"
if diff grep_output7.txt my_grep_output7.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 8: Контекст до (-B 2) ==="
echo "GNU grep -B 2 'date':"
grep -B 2 'date' test_input.txt > grep_output8.txt
cat grep_output8.txt

echo -e "\nMy grep -B 2 'date':"
./my_grep -B 2 'date' test_input.txt > my_grep_output8.txt
cat my_grep_output8.txt

echo -e "\nСравнение:"
if diff grep_output8.txt my_grep_output8.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 9: Полный контекст (-C 1) ==="
echo "GNU grep -C 1 'grape':"
grep -C 1 'grape' test_input.txt > grep_output9.txt
cat grep_output9.txt

echo -e "\nMy grep -C 1 'grape':"
./my_grep -C 1 'grape' test_input.txt > my_grep_output9.txt
cat my_grep_output9.txt

echo -e "\nСравнение:"
if diff grep_output9.txt my_grep_output9.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 10: Комбинированные флаги (-i -n -C 1) ==="
echo "GNU grep -i -n -C 1 'apple':"
grep -i -n -C 1 'apple' test_input.txt > grep_output10.txt
cat grep_output10.txt

echo -e "\nMy grep -i -n -C 1 'apple':"
./my_grep -i -n -C 1 'apple' test_input.txt > my_grep_output10.txt
cat my_grep_output10.txt

echo -e "\nСравнение:"
if diff grep_output10.txt my_grep_output10.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 11: Регулярное выражение ==="
echo "GNU grep 'line [0-9]':"
grep 'line [0-9]' test_input.txt > grep_output11.txt
cat grep_output11.txt | head -5

echo -e "\nMy grep 'line [0-9]':"
./my_grep 'line [0-9]' test_input.txt > my_grep_output11.txt
cat my_grep_output11.txt | head -5

echo -e "\nСравнение:"
if diff grep_output11.txt my_grep_output11.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 12: Из stdin ==="
echo "GNU grep из stdin 'banana':"
cat test_input.txt | grep 'banana' > grep_output12.txt
cat grep_output12.txt

echo -e "\nMy grep из stdin 'banana':"
cat test_input.txt | go run *.go 'banana' > my_grep_output12.txt
cat my_grep_output12.txt

echo -e "\nСравнение:"
if diff grep_output12.txt my_grep_output12.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Сводка ==="
echo "Созданные файлы:"
ls -la *output*.txt test_input.txt

echo -e "\nДля детального сравнения используйте:"
echo "  diff -u grep_outputX.txt my_grep_outputX.txt"
echo "  meld grep_outputX.txt my_grep_outputX.txt"