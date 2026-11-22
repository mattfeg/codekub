# Dígitos de um número

Separe um número de três dígitos em centena, dezena e unidade, e some esses
dígitos, usando só operações aritméticas.

## Entrada

Um inteiro `N` com exatamente três dígitos (100 ≤ N ≤ 999).

## Saída

Imprima quatro linhas:

```
centena: <dígito das centenas>
dezena: <dígito das dezenas>
unidade: <dígito das unidades>
soma: <soma dos três dígitos>
```

## Exemplos

Entrada:

```
407
```

Saída:

```
centena: 4
dezena: 0
unidade: 7
soma: 11
```

Entrada:

```
999
```

Saída:

```
centena: 9
dezena: 9
unidade: 9
soma: 27
```

## Dica

O resto da divisão por 10 é sempre o último dígito. E dividir por 10 (divisão
inteira) "apaga" o último dígito.
