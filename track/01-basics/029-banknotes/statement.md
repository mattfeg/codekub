# Troco em notas

Um caixa eletrônico precisa entregar um valor usando a menor quantidade
possível de cédulas e moedas. Os valores disponíveis são R$ 200, 100, 50, 20,
10, 5, 2 e 1.

Com esses valores, a menor quantidade sai de um jeito simples: use o máximo
possível do maior valor, depois o máximo do próximo com o que sobrou, e assim
por diante.

## Entrada

Um inteiro `V` (0 ≤ V ≤ 10⁶), o valor em reais.

## Saída

Imprima oito linhas, do maior para o menor valor, no formato
`R$ <valor>: <quantidade>`.

## Exemplos

Entrada:

```
388
```

Saída:

```
R$ 200: 1
R$ 100: 1
R$ 50: 1
R$ 20: 1
R$ 10: 1
R$ 5: 1
R$ 2: 1
R$ 1: 1
```

Entrada:

```
76
```

Saída:

```
R$ 200: 0
R$ 100: 0
R$ 50: 1
R$ 20: 1
R$ 10: 0
R$ 5: 1
R$ 2: 0
R$ 1: 1
```

## Dica

Depois de calcular quantas cédulas de um valor cabem (divisão inteira), o
valor que sobra para as próximas é o **resto** dessa divisão.
