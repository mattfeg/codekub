# Segundos em horas

Um cronômetro marca o tempo só em segundos. Mostre esse tempo no formato de
relógio `H:MM:SS`.

## Entrada

Um inteiro `N` (0 ≤ N ≤ 10⁷), a quantidade de segundos.

## Saída

Imprima `H:MM:SS`, onde:

- `H` é o número de horas, sem zeros à esquerda (pode passar de 24);
- `MM` são os minutos e `SS` os segundos, **sempre com dois dígitos**.

## Exemplos

Entrada:

```
3725
```

Saída:

```
1:02:05
```

Entrada:

```
0
```

Saída:

```
0:00:00
```

## Dica

Uma hora tem 3600 segundos e um minuto tem 60. Depois de tirar as horas
inteiras, o **resto** é o que sobra para minutos e segundos.
