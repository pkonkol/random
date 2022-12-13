import functools

lines = list(map(str.splitlines, open("input").read().strip().split("\n\n")))

# return >0 if left is greater?
def check(left, right):
    # if both are int
    if type(left) == int:
        if type(right) == int:
            return left - right
        else:
            return check([left], right)
    # if second is int
    if type(right) == int:
        return check(left, [right])
    # if both are lists
    for i, (l, r) in enumerate(zip(left, right)):
        cmp = check(l, r)
        if cmp:
            return cmp

    # if zipped checks are equal but there may be more
    return len(left) - len(right)

sum = 0
for i, (left, right) in enumerate(lines):
    x = check(eval(left), eval(right))
    print(i, " ", eval(left), " ", eval(right), " ", x)
    if x < 0:
        sum += i + 1

print(f"sum for task 1 is {sum}")

lines = list(map(str.splitlines, open("input").read().strip().split("\n\n")))
merged = [eval(l) for (l, _) in lines] + [eval(r) for (_, r) in lines] +  [[[2]]] + [[[6]]]
merged.sort(key=functools.cmp_to_key(lambda x, y: check(x, y)))
# for i in merged:
#     print(i)
print((merged.index([[2]])+1) * (merged.index([[6]])+1))