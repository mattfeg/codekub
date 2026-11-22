# Consumo de combustível

Depois de uma viagem, Laura quer saber quantos quilômetros o carro fez por
litro e quanto gastou, em média, a cada quilômetro.

## Entrada

Três linhas:

1. a distância percorrida, em km: um inteiro `D` (1 ≤ D ≤ 100000);
2. os litros gastos: um real `L` (0.01 ≤ L ≤ 10000), com duas casas;
3. o preço do litro: um real `P` (0.01 ≤ P ≤ 20), com duas casas.

## Saída

Imprima duas linhas:

```
Consumo: <D / L, com três casas decimais> km/l
Custo por km: R$ <L × P / D, com duas casas decimais>
```

## Exemplos

Entrada:

```
500
40.00
6.00
```

Saída:

```
Consumo: 12.500 km/l
Custo por km: R$ 0.48
```

Entrada:

```
123
9.87
5.49
```

Saída:

```
Consumo: 12.462 km/l
Custo por km: R$ 0.44
```
