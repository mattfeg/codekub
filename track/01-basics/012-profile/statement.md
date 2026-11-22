# Ficha cadastral

Uma escola de programação está montando fichas de cadastro dos alunos. Cada
ficha mostra os dados lidos e uma previsão simpática para o futuro.

## Entrada

Três linhas:

1. o nome completo (1 a 60 caracteres, pode ter espaços);
2. a idade, um inteiro `I` (0 ≤ I ≤ 120);
3. a cidade (1 a 40 caracteres, pode ter espaços).

## Saída

Imprima quatro linhas:

```
Nome: <nome>
Idade: <idade> anos
Cidade: <cidade>
Em 10 anos, <nome> terá <idade + 10> anos.
```

## Exemplos

Entrada:

```
Ana Souza
28
Recife
```

Saída:

```
Nome: Ana Souza
Idade: 28 anos
Cidade: Recife
Em 10 anos, Ana Souza terá 38 anos.
```

Entrada:

```
Bruno
0
São José dos Campos
```

Saída:

```
Nome: Bruno
Idade: 0 anos
Cidade: São José dos Campos
Em 10 anos, Bruno terá 10 anos.
```

## Observação

A idade chega como texto. Para somar 10, você precisa convertê-la para um
número inteiro.
