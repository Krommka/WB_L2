#!/bin/bash

# Тестовый скрипт для сравнения my_cut с GNU cut

# Создаем тестовый файл с табуляцией
cat > test_input.txt << EOF
apple	5	Jan	fruit
banana	2	Mar	fruit
cherry	1	Feb	fruit
date	3	Dec	fruit
elderberry	4	Nov	berry
fig	6	Jul	fruit
grape	7	Apr	fruit
honeydew	8	Jun	melon
kiwi	9	May	fruit
lemon	10	Aug	fruit
mango	11	Sep	fruit
nectarine	12	Oct	fruit
single_column
EOF

echo "=== Тестовый файл создан ==="
cat test_input.txt | sed 's/\t/\\t/g'

echo -e "\n=== Тест 1: Простой выбор поля (-f 1) ==="
echo "GNU cut -f 1:"
cut -f 1 test_input.txt > cut_output1.txt
cat cut_output1.txt

echo -e "\nMy cut -f 1:"
./my_cut -f 1 test_input.txt > my_cut_output1.txt
cat my_cut_output1.txt

echo -e "\nСравнение:"
if diff cut_output1.txt my_cut_output1.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 2: Несколько полей (-f 1,3) ==="
echo "GNU cut -f 1,3:"
cut -f 1,3 test_input.txt > cut_output2.txt
cat cut_output2.txt

echo -e "\nMy cut -f 1,3:"
./my_cut -f 1,3 test_input.txt > my_cut_output2.txt
cat my_cut_output2.txt

echo -e "\nСравнение:"
if diff cut_output2.txt my_cut_output2.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 3: Диапазон полей (-f 2-4) ==="
echo "GNU cut -f 2-4:"
cut -f 2-4 test_input.txt > cut_output3.txt
cat cut_output3.txt

echo -e "\nMy cut -f 2-4:"
./my_cut -f 2-4 test_input.txt > my_cut_output3.txt
cat my_cut_output3.txt

echo -e "\nСравнение:"
if diff cut_output3.txt my_cut_output3.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 4: Комбинированные поля и диапазоны (-f 1,3-4) ==="
echo "GNU cut -f 1,3-4:"
cut -f 1,3-4 test_input.txt > cut_output4.txt
cat cut_output4.txt

echo -e "\nMy cut -f 1,3-4:"
./my_cut -f 1,3-4 test_input.txt > my_cut_output4.txt
cat my_cut_output4.txt

echo -e "\nСравнение:"
if diff cut_output4.txt my_cut_output4.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 5: Только строки с разделителем (-s -f 1) ==="
echo "GNU cut -s -f 1:"
cut -s -f 1 test_input.txt > cut_output5.txt
cat cut_output5.txt

echo -e "\nMy cut -s -f 1:"
./my_cut -s -f 1 test_input.txt > my_cut_output5.txt
cat my_cut_output5.txt

echo -e "\nСравнение:"
if diff cut_output5.txt my_cut_output5.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

# Создаем CSV файл для тестов с другим разделителем
cat > test_input.csv << EOF
apple,5,Jan,fruit
banana,2,Mar,fruit
cherry,1,Feb,fruit
date,3,Dec,fruit
single_column
EOF

echo -e "\n=== Тест 6: Другой разделитель (-d ',' -f 1,3) ==="
echo "GNU cut -d ',' -f 1,3:"
cut -d ',' -f 1,3 test_input.csv > cut_output6.txt
cat cut_output6.txt

echo -e "\nMy cut -d ',' -f 1,3:"
./my_cut -d ',' -f 1,3 test_input.csv > my_cut_output6.txt
cat my_cut_output6.txt

echo -e "\nСравнение:"
if diff cut_output6.txt my_cut_output6.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 7: Комбинированные флаги (-s -d ',' -f 2) ==="
echo "GNU cut -s -d ',' -f 2:"
cut -s -d ',' -f 2 test_input.csv > cut_output7.txt
cat cut_output7.txt

echo -e "\nMy cut -s -d ',' -f 2:"
./my_cut -s -d ',' -f 2 test_input.csv > my_cut_output7.txt
cat my_cut_output7.txt

echo -e "\nСравнение:"
if diff cut_output7.txt my_cut_output7.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 8: Поля за границами (-f 1,5,10) ==="
echo "GNU cut -f 1,5,10:"
cut -f 1,5,10 test_input.txt > cut_output8.txt
cat cut_output8.txt

echo -e "\nMy cut -f 1,5,10:"
./my_cut -f 1,5,10 test_input.txt > my_cut_output8.txt
cat my_cut_output8.txt

echo -e "\nСравнение:"
if diff cut_output8.txt my_cut_output8.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 9: Из stdin ==="
echo "GNU cut из stdin -f 2:"
cat test_input.txt | cut -f 2 > cut_output9.txt
cat cut_output9.txt

echo -e "\nMy cut из stdin -f 2:"
cat test_input.txt | ./my_cut -f 2 > my_cut_output9.txt
cat my_cut_output9.txt

echo -e "\nСравнение:"
if diff cut_output9.txt my_cut_output9.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Тест 10: Сложный диапазон (-f 1-2,4) ==="
echo "GNU cut -f 1-2,4:"
cut -f 1-2,4 test_input.txt > cut_output10.txt
cat cut_output10.txt

echo -e "\nMy cut -f 1-2,4:"
./my_cut -f 1-2,4 test_input.txt > my_cut_output10.txt
cat my_cut_output10.txt

echo -e "\nСравнение:"
if diff cut_output10.txt my_cut_output10.txt; then
    echo "✓ Результаты идентичны"
else
    echo "✗ Результаты различаются"
fi

echo -e "\n=== Сводка ==="
echo "Созданные файлы:"
ls -la *output*.txt test_input.txt test_input.csv

echo -e "\nДля детального сравнения используйте:"
echo "  diff -u cut_outputX.txt my_cut_outputX.txt"
echo "  meld cut_outputX.txt my_cut_outputX.txt"