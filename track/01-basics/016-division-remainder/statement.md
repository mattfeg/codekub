# Divisão inteira e resto

Na divisão que aprendemos na escola, dividir 17 por 5 dá 3 e sobra 2. Essas
duas partes, o **quociente** e o **resto**, aparecem o tempo todo em
programação.

## Entrada

Uma linha com dois inteiros `A` e `B` (0 ≤ A ≤ 10⁹, 1 ≤ B ≤ 10⁹).

## Saída

Imprima três linhas:

```
quociente: <A dividido por B, sem a parte fracionária>
resto: <o que sobra da divisão>
<A> = <quociente> x <B> + <resto>
```

## Exemplos

Entrada:

```
17 5
```

Saída:

```
quociente: 3
resto: 2
17 = 3 x 5 + 2
```

Entrada:

```
4 9
```

Saída:

```
quociente: 0
resto: 4
4 = 0 x 9 + 4
```

## Dica

A maioria das linguagens tem um operador próprio para o resto, geralmente `%`.
