# Piso, teto e arredondamento

Existem várias formas de transformar um número real em inteiro, e cada uma dá
um resultado diferente com números negativos:

- **piso**: o maior inteiro menor ou igual a `X`;
- **teto**: o menor inteiro maior ou igual a `X`;
- **truncado**: `X` sem a parte fracionária (corta em direção ao zero);
- **arredondado**: o inteiro mais próximo de `X`. Quando `X` está exatamente
  no meio (como `2.5`), arredonde para **longe do zero**: `2.5` vira `3` e
  `-2.5` vira `-3`.

## Entrada

Um número real `X` (−1000 ≤ X ≤ 1000), com até três casas decimais.

## Saída

Imprima quatro linhas com os valores como **inteiros**, sem casas decimais:

```
piso: <piso>
teto: <teto>
truncado: <truncado>
arredondado: <arredondado>
```

## Exemplos

Entrada:

```
2.5
```

Saída:

```
piso: 2
teto: 3
truncado: 2
arredondado: 3
```

Entrada:

```
-2.5
```

Saída:

```
piso: -3
teto: -2
truncado: -2
arredondado: -3
```

## Observação

Nem toda linguagem arredonda `.5` para longe do zero. O `round` do Python, por
exemplo, arredonda para o par mais próximo (`round(2.5)` dá `2`). Confira como
a sua funciona. E imprima `0`, não `-0`, quando o resultado for zero.
