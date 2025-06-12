// 基于 Vec 实现的泛型栈数据结构

/// 一个基于 Vec 实现的泛型栈
/// 
/// # 类型参数
/// 
/// * `T` - 栈中存储的元素类型
/// 
/// # 示例
/// 
/// ```
/// use rust_playground::stack::Stack;
/// 
/// let mut stack = Stack::new();
/// stack.push(1);
/// stack.push(2);
/// stack.push(3);
/// 
/// assert_eq!(stack.pop(), Some(3));
/// assert_eq!(stack.peek(), Some(&2));
/// assert_eq!(stack.size(), 2);
/// ```
pub struct Stack<T> {
    elements: Vec<T>,
}

impl<T> Stack<T> {
    /// 创建一个新的空栈
    pub fn new() -> Self {
        Stack {
            elements: Vec::new(),
        }
    }

    /// 使用指定容量创建一个新的空栈
    pub fn with_capacity(capacity: usize) -> Self {
        Stack {
            elements: Vec::with_capacity(capacity),
        }
    }

    /// 将元素压入栈顶
    pub fn push(&mut self, element: T) {
        self.elements.push(element);
    }

    /// 弹出栈顶元素
    /// 
    /// 如果栈为空，返回 None
    pub fn pop(&mut self) -> Option<T> {
        self.elements.pop()
    }

    /// 查看栈顶元素但不移除
    /// 
    /// 如果栈为空，返回 None
    pub fn peek(&self) -> Option<&T> {
        self.elements.last()
    }

    /// 返回栈中元素的数量
    pub fn size(&self) -> usize {
        self.elements.len()
    }

    /// 检查栈是否为空
    pub fn is_empty(&self) -> bool {
        self.elements.is_empty()
    }

    /// 清空栈
    pub fn clear(&mut self) {
        self.elements.clear();
    }

    /// 返回栈中所有元素的迭代器
    pub fn iter(&self) -> impl Iterator<Item = &T> {
        self.elements.iter()
    }
}

// 为 Stack 实现 Default trait
impl<T> Default for Stack<T> {
    fn default() -> Self {
        Self::new()
    }
}

// 为 Stack 实现 IntoIterator trait，允许使用 for 循环遍历栈
impl<T> IntoIterator for Stack<T> {
    type Item = T;
    type IntoIter = std::vec::IntoIter<Self::Item>;

    fn into_iter(self) -> Self::IntoIter {
        self.elements.into_iter()
    }
}

// 为 Stack 实现 FromIterator trait，允许从迭代器创建栈
impl<T> FromIterator<T> for Stack<T> {
    fn from_iter<I: IntoIterator<Item = T>>(iter: I) -> Self {
        let mut stack = Stack::new();
        for item in iter {
            stack.push(item);
        }
        stack
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_push_pop() {
        let mut stack = Stack::new();
        stack.push(1);
        stack.push(2);
        stack.push(3);

        assert_eq!(stack.pop(), Some(3));
        assert_eq!(stack.pop(), Some(2));
        assert_eq!(stack.pop(), Some(1));
        assert_eq!(stack.pop(), None);
    }

    #[test]
    fn test_peek() {
        let mut stack = Stack::new();
        assert_eq!(stack.peek(), None);

        stack.push(1);
        assert_eq!(stack.peek(), Some(&1));

        stack.push(2);
        assert_eq!(stack.peek(), Some(&2));

        stack.pop();
        assert_eq!(stack.peek(), Some(&1));
    }

    #[test]
    fn test_size_and_empty() {
        let mut stack = Stack::new();
        assert_eq!(stack.size(), 0);
        assert!(stack.is_empty());

        stack.push(1);
        assert_eq!(stack.size(), 1);
        assert!(!stack.is_empty());

        stack.push(2);
        assert_eq!(stack.size(), 2);

        stack.pop();
        assert_eq!(stack.size(), 1);

        stack.pop();
        assert_eq!(stack.size(), 0);
        assert!(stack.is_empty());
    }

    #[test]
    fn test_clear() {
        let mut stack = Stack::new();
        stack.push(1);
        stack.push(2);
        stack.push(3);

        stack.clear();
        assert_eq!(stack.size(), 0);
        assert!(stack.is_empty());
        assert_eq!(stack.pop(), None);
    }

    #[test]
    fn test_iter() {
        let mut stack = Stack::new();
        stack.push(1);
        stack.push(2);
        stack.push(3);

        let mut iter = stack.iter();
        assert_eq!(iter.next(), Some(&1));
        assert_eq!(iter.next(), Some(&2));
        assert_eq!(iter.next(), Some(&3));
        assert_eq!(iter.next(), None);
    }

    #[test]
    fn test_into_iter() {
        let mut stack = Stack::new();
        stack.push(1);
        stack.push(2);
        stack.push(3);

        let mut values = Vec::new();
        for value in stack {
            values.push(value);
        }

        assert_eq!(values, vec![1, 2, 3]);
    }

    #[test]
    fn test_from_iter() {
        let values = vec![1, 2, 3];
        let stack: Stack<i32> = values.into_iter().collect();

        assert_eq!(stack.size(), 3);
        assert_eq!(stack.peek(), Some(&3));
    }
}
