# Horário de chegada

Um ônibus sai em um certo horário e a viagem dura alguns minutos, às vezes
vários dias. Descubra em que horário ele chega e quantos dias depois da
partida.

## Entrada

Duas linhas:

1. a hora `H` e o minuto `M` da partida (0 ≤ H ≤ 23, 0 ≤ M ≤ 59);
2. a duração da viagem em minutos, `D` (0 ≤ D ≤ 10⁶).

## Saída

Imprima duas linhas:

```
chegada: HH:MM
dias depois: <quantos dias após o dia da partida>
```

A hora e o minuto da chegada têm **sempre dois dígitos**.

## Exemplos

Entrada:

```
22 30
525
```

Saída:

```
chegada: 07:15
dias depois: 1
```

Entrada:

```
8 0
90
```

Saída:

```
chegada: 09:30
dias depois: 0
```

## Dica

Converta tudo para minutos contados desde a meia-noite do dia da partida. Um
dia tem 1440 minutos.
