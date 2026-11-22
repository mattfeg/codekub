# Expressões e precedência

Em matemática e em programação, a multiplicação e a divisão são feitas antes
da soma e da subtração, e operações de mesma prioridade são feitas da
esquerda para a direita. Os parênteses mudam essa ordem.

Dados quatro inteiros `a`, `b`, `c` e `d`, calcule as expressões abaixo. Todas
as divisões são **inteiras** (descartam a parte fracionária) e `%` é o resto
da divisão.

- `X = a + b * c - d`
- `Y = (a + b) * (c - d)`
- `Z = (a * b + c) % d`
- `W = a + b / c * d`

## Entrada

Uma linha com quatro inteiros `a`, `b`, `c` e `d` (1 ≤ a, b, c, d ≤ 1000).

## Saída

Imprima quatro linhas: `X = ...`, `Y = ...`, `Z = ...` e `W = ...`.

## Exemplos

Entrada:

```
2 3 4 5
```

Saída:

```
X = 9
Y = -5
Z = 0
W = 2
```

Entrada:

```
10 7 3 4
```

Saída:

```
X = 27
Y = -17
Z = 1
W = 18
```

## Observação

Em `W`, `b / c * d` é calculado como `(b / c) * d`, da esquerda para a direita.
No primeiro exemplo, `3 / 4` dá `0`, e por isso `W = 2 + 0 * 5 = 2`.
