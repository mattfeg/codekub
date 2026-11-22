# Conta do restaurante

Um grupo de amigos vai dividir a conta do jantar igualmente, já com a taxa de
serviço incluída.

## Entrada

Uma linha com três valores:

- o valor da conta `C`, um real com duas casas (0.01 ≤ C ≤ 10000);
- a taxa de serviço `S`, em porcentagem, um inteiro (0 ≤ S ≤ 20);
- o número de pessoas `N`, um inteiro (1 ≤ N ≤ 50).

## Saída

Imprima três linhas, com valores de duas casas decimais:

```
Serviço: R$ <C × S / 100>
Total: R$ <C mais o serviço>
Por pessoa: R$ <total dividido por N>
```

## Exemplos

Entrada:

```
150.00 10 4
```

Saída:

```
Serviço: R$ 15.00
Total: R$ 165.00
Por pessoa: R$ 41.25
```

Entrada:

```
89.90 0 3
```

Saída:

```
Serviço: R$ 0.00
Total: R$ 89.90
Por pessoa: R$ 29.97
```
