# Peças de dominó

> **Chefão** · nota 800 no Codeforces

Você tem um tabuleiro retangular de `M × N` casas e quantas peças de dominó
quiser. Cada peça ocupa exatamente **duas casas vizinhas** (`2 × 1`) e pode ser
girada.

Coloque o maior número possível de peças no tabuleiro, respeitando as regras:

1. cada peça cobre completamente duas casas;
2. duas peças não podem se sobrepor;
3. cada peça fica inteira dentro do tabuleiro (pode encostar na borda).

Qual é o número máximo de peças?

## Entrada

Uma linha com dois inteiros `M` e `N`, as dimensões do tabuleiro em casas
(1 ≤ M ≤ N ≤ 16).

## Saída

Imprima um único número: a quantidade máxima de peças.

## Exemplos

Entrada:

```
2 4
```

Saída:

```
4
```

Entrada:

```
3 3
```

Saída:

```
4
```

## Dica

Não é preciso simular a colocação das peças. Pense em quantas casas o
tabuleiro tem e em quantas casas cada peça ocupa.

---

*Adaptado de [Codeforces 50A — Domino piling](https://codeforces.com/problemset/problem/50/A),
via dataset open-r1/codeforces (CC BY 4.0). Traduzido para o português.*
