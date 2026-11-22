# Latas de tinta

Para pintar as quatro paredes de um quarto, é preciso saber a área total e
quantas latas de tinta comprar. Ignore portas e janelas.

A área das quatro paredes é `2 × (comprimento + largura) × altura`, e ela
precisa ser multiplicada pelo número de demãos (camadas de tinta). Cada lata
cobre uma certa área, e só dá para comprar latas inteiras.

## Entrada

Duas linhas:

1. o comprimento `C`, a largura `L` e a altura `A` do quarto, em metros
   (inteiros, 1 ≤ C, L, A ≤ 100);
2. o número de demãos `D` (1 ≤ D ≤ 5) e quantos m² cada lata cobre, `R`
   (1 ≤ R ≤ 1000).

## Saída

Imprima duas linhas:

```
Área total: <área a pintar, contando as demãos> m²
Latas: <quantidade mínima de latas>
```

## Exemplos

Entrada:

```
4 3 3
2 36
```

Saída:

```
Área total: 84 m²
Latas: 3
```

Entrada:

```
5 5 2
1 40
```

Saída:

```
Área total: 40 m²
Latas: 1
```

## Dica

84 m² com latas de 36 m² dá 2,33…, então são 3 latas: é preciso arredondar
**para cima**. Dá para fazer isso só com divisão inteira, sem números reais:
pense no que acontece com `(área + R − 1) / R`.
