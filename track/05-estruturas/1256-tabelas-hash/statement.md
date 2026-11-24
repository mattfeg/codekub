# Tabelas Hash

O enunciado está no beecrowd: aperte `o` para abri-lo no navegador.

<https://resources.beecrowd.com/repository/UOJ_1256.html>

Os testes são só os exemplos do enunciado: passar neles não garante que a
solução esteja certa no beecrowd.

## Exemplos

Entrada:

```
2
13 9
44 45 49 70 27 73 92 97 95
7 8
35 12 2 17 19 51 88 86
```

Saída:

```
0 -> \
1 -> 27 -> 92 -> \
2 -> \
3 -> \
4 -> 95 -> \
5 -> 44 -> 70 -> \
6 -> 45 -> 97 -> \
7 -> \
8 -> 73 -> \
9 -> \
10 -> 49 -> \
11 -> \
12 -> \

0 -> 35 -> \
1 -> \
2 -> 2 -> 51 -> 86 -> \
3 -> 17 -> \
4 -> 88 -> \
5 -> 12 -> 19 -> \
6 -> \
```
