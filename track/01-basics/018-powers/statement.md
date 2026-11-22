# Potências

Dado um inteiro `N`, mostre as potências de `N` com expoentes de 2 a 4.

## Entrada

Um inteiro `N` (0 ≤ N ≤ 1000).

## Saída

Imprima três linhas:

```
N^2 = <N elevado a 2>
N^3 = <N elevado a 3>
N^4 = <N elevado a 4>
```

## Exemplos

Entrada:

```
3
```

Saída:

```
N^2 = 9
N^3 = 27
N^4 = 81
```

Entrada:

```
1000
```

Saída:

```
N^2 = 1000000
N^3 = 1000000000
N^4 = 1000000000000
```

## Observação

1000⁴ = 10¹² não cabe em um inteiro de 32 bits. Em linguagens como C, C++, Java
e Go, use um tipo de 64 bits (`long long`, `long`, `int64`). E fuja de funções
de potência que devolvem número real: elas podem imprimir `1e+12` ou casas
decimais.
