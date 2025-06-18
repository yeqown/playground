/*
 * 迭代器是一种数据结构，它允许我们遍历一个集合中的元素，而不需要知道集合的内部结构。
 * 迭代器是一个对象，它实现了 next() 方法，每次调用 next() 方法，都会返回集合中的下一个元素。
 *
 * 1. 迭代器是惰性的，只有在调用 next() 方法时，才会返回元素。
 * 2. 要实现迭代器，需要实现 Iterator Trait, 该 Trait 定义了 next() 方法。
 * 3. IntoIterator Trait 定义了 into_iter() 方法，用于将将数据结构转换为迭代器。
 *
 * xxx.into_iter() 转移所有权, 不能再使用
 * xxx.iter()      不可变借用，不能修改元素
 * xxx.iter_mut()  可变借用，可以修改元素
 */
use std::iter::IntoIterator;

fn main() {
    let v = vec![1, 2, 3];

    let mut iter = v.into_iter();

    // println!("{:?}", v); // v 已经被转移所有权，不能再使用

    assert_eq!(iter.next(), Some(1));
    assert_eq!(iter.next(), Some(2));
    assert_eq!(iter.next(), Some(3));
    assert_eq!(iter.next(), None);

    let v2 = vec![1, 2, 3];
    // v2.iter().map(|x| x + 1); // Unused Map<Iter<i32>, fn(&i32) -> i32> that must be used
    let plus_one = v2.iter().map(|x| x + 1).collect::<Vec<i32>>();
    println!("{:?}", plus_one);

    println!("v2 can be used again:");
    v2.iter().for_each(|x| {
        println!("{}", x);
    });
}