# Dia da semana

Os dias da semana são numerados assim:

| Número | Dia           |
|--------|---------------|
| 0      | domingo       |
| 1      | segunda-feira |
| 2      | terça-feira   |
| 3      | quarta-feira  |
| 4      | quinta-feira  |
| 5      | sexta-feira   |
| 6      | sábado        |

Sabendo o dia de hoje, descubra que dia da semana será daqui a `N` dias. Se
`N` for negativo, a pergunta é que dia foi `|N|` dias atrás.

## Entrada

Uma linha com dois inteiros: o dia de hoje `D` (0 ≤ D ≤ 6) e `N`
(−10⁹ ≤ N ≤ 10⁹).

## Saída

Imprima `dia: ` seguido do número do dia da semana, de 0 a 6.

## Exemplos

Entrada:

```
3 10
```

Saída:

```
dia: 6
```

Entrada:

```
1 -3
```

Saída:

```
dia: 5
```

## Observação

No segundo exemplo, 3 dias antes de segunda-feira é sexta-feira (5). Cuidado:
em linguagens como C, C++, Java, Go e JavaScript, o resto de um número negativo
também é negativo (`-2 % 7` dá `-2`). Em Ruby e Python, dá `5`. Sua resposta
precisa estar sempre entre 0 e 6.
