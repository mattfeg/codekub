# Média de três notas

Um professor dá notas inteiras de 0 a 100 e, no boletim, mostra só a parte
inteira da média (sem arredondar para cima).

## Entrada

Uma linha com três inteiros `A`, `B` e `C` (0 ≤ A, B, C ≤ 100).

## Saída

Imprima `Média: M`, onde `M` é a parte inteira de `(A + B + C) / 3`.

## Exemplos

Entrada:

```
70 80 90
```

Saída:

```
Média: 80
```

Entrada:

```
100 100 99
```

Saída:

```
Média: 99
```

## Observação

No segundo exemplo, a média exata é 99,666…, mas o boletim mostra só a parte
inteira: `99`. Repare que a divisão de inteiros em muitas linguagens já
descarta a parte fracionária.
