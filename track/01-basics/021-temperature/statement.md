# Conversão de temperatura

Converta uma temperatura de graus Celsius para Fahrenheit e para Kelvin:

- `F = C × 9 / 5 + 32`
- `K = C + 273.15`

## Entrada

Um número real `C` (−273.15 ≤ C ≤ 10000), com até duas casas decimais.

## Saída

Imprima duas linhas, cada valor com **duas** casas decimais:

```
Fahrenheit: <F>
Kelvin: <K>
```

## Exemplos

Entrada:

```
37.00
```

Saída:

```
Fahrenheit: 98.60
Kelvin: 310.15
```

Entrada:

```
-40.00
```

Saída:

```
Fahrenheit: -40.00
Kelvin: 233.15
```

## Observação

Cuidado com a divisão: em algumas linguagens, `9 / 5` entre inteiros dá `1`, e
não `1.8`. Garanta que a conta seja feita com números reais.
