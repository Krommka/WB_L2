#!/bin/bash

# Скрипт для сравнения my_sort и GNU sort

# Создаем тестовый файл
cat > test_input.txt << EOF
apple	5	Jan
banana	2K	Mar
cherry	1	Feb
date	500M	Dec
elderberry	3	Nov
fig	2G	Jul
grape	100	Apr
honeydew	1K	Jun
kiwi	50	May
lemon	200	Aug
mango	1.5K	Sep
nectarine	300	Oct
apple	5	Jan
banana	2K	Mar
EOF

echo "=== Тестовый файл создан ==="
cat test_input.txt

echo -e "\n=== Тест 1: Простая сортировка ==="
echo "GNU sort:"
sort test_input.txt > sort_output1.txt
cat sort_output1.txt

echo -e "\nMy sort:"
go run my_sort.go test_input.txt > my_sort_output1.txt
cat my_sort_output1.txt

echo -e "\nСравнение:"
if diff sort_output1.txt my_sort_output1.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 2: Числовая сортировка по колонке 2 (-k 2 -n) ==="
echo "GNU sort:"
sort -k 2,2 -n test_input.txt > sort_output2.txt
cat sort_output2.txt

echo -e "\nMy sort:"
go run my_sort.go -k 2 -n test_input.txt > my_sort_output2.txt
cat my_sort_output2.txt

echo -e "\nСравнение:"
if diff sort_output2.txt my_sort_output2.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 3: Обратная сортировка (-r) ==="
echo "GNU sort:"
sort -r test_input.txt > sort_output3.txt
cat sort_output3.txt

echo -e "\nMy sort:"
go run my_sort.go -r test_input.txt > my_sort_output3.txt
cat my_sort_output3.txt

echo -e "\nСравнение:"
if diff sort_output3.txt my_sort_output3.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 4: Уникальные строки (-u) ==="
echo "GNU sort:"
sort -u test_input.txt > sort_output4.txt
cat sort_output4.txt

echo -e "\nMy sort:"
go run my_sort.go -u test_input.txt > my_sort_output4.txt
cat my_sort_output4.txt

echo -e "\nСравнение:"
if diff sort_output4.txt my_sort_output4.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 5: Комбинированные флаги (-k 2 -nr) ==="
echo "GNU sort:"
sort -k 2,2 -nr test_input.txt > sort_output5.txt
cat sort_output5.txt

echo -e "\nMy sort:"
go run my_sort.go -k 2 -nr test_input.txt > my_sort_output5.txt
cat my_sort_output5.txt

echo -e "\nСравнение:"
if diff sort_output5.txt my_sort_output5.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 6: Сортировка по месяцам (-k 3 -M) ==="
echo "GNU sort:"
sort -k 3,3 -M test_input.txt > sort_output6.txt
cat sort_output6.txt

echo -e "\nMy sort:"
go run my_sort.go -k 3 -M test_input.txt > my_sort_output6.txt
cat my_sort_output6.txt

echo -e "\nСравнение:"
if diff sort_output6.txt my_sort_output6.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 7: Человеко-читаемые числа (-k 2 -h) ==="
echo "GNU sort:"
sort -k 2,2 -h test_input.txt > sort_output7.txt
cat sort_output7.txt

echo -e "\nMy sort:"
go run my_sort.go -k 2 -h test_input.txt > my_sort_output7.txt
cat my_sort_output7.txt

echo -e "\nСравнение:"
if diff sort_output7.txt my_sort_output7.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 8: Проверка сортировки (-c) ==="
echo "GNU sort:"
sort -c test_input.txt 2>&1 || true

echo -e "\nMy sort:"
go run my_sort.go -c test_input.txt 2>&1 || true

# Создаем отсортированный файл для проверки
sort test_input.txt > sorted_test.txt
echo -e "\nGNU sort на отсортированном файле:"
sort -c sorted_test.txt 2>&1 || true

echo -e "\nMy sort на отсортированном файле:"
go run my_sort.go -c sorted_test.txt 2>&1 || true

echo -e "\n=== Сводка ==="
echo "Созданные файлы:"
ls -la *output*.txt sorted_test.txt test_input.txt

echo -e "\nДля детального сравнения используйте:"
echo "  diff -u sort_outputX.txt my_sort_outputX.txt"
echo "  meld sort_outputX.txt my_sort_outputX.txt"

# Очистка (раскомментируйте если нужно)
# rm -f test_input.txt sorted_test.txt *output*.txt