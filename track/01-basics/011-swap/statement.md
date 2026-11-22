# Troca de valores

Duas variáveis, `A` e `B`, guardam um número cada. Seu programa deve trocar o
conteúdo delas: o que estava em `A` vai para `B`, e o que estava em `B` vai
para `A`.

## Entrada

Uma linha com dois inteiros `A` e `B` (−10⁹ ≤ A, B ≤ 10⁹).

## Saída

Imprima duas linhas:

```
antes: A = <valor de A>, B = <valor de B>
depois: A = <novo valor de A>, B = <novo valor de B>
```

## Exemplos

Entrada:

```
3 7
```

Saída:

```
antes: A = 3, B = 7
depois: A = 7, B = 3
```

Entrada:

```
-1 -1
```

Saída:

```
antes: A = -1, B = -1
depois: A = -1, B = -1
```

## Dica

Se você fizer `A = B` e depois `B = A`, as duas variáveis terminam com o
mesmo valor. Pense em como guardar um dos valores antes de sobrescrevê-lo.
