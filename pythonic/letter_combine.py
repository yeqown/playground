class Solution:
    m = {
        "2": ["a", "b", "c"],
        "3": ["d", "e", "f"],
        "4": ["g", "h", "i"],
        "5": ["j", "k", "l"],
        "6": ["m", "n", "o"],
        "7": ["p", "q", "r", "s"],
        "8": ["t", "u", "v"],
        "9": ["w", "x", "y", "z"]
    }

    def multiply(self, a, b):
        res = []
        for ac in a:
            for bc in b:
                res.append(ac+bc)
        return res

    def letterCombinations(self, digits):
        """
        :type digits: str
        :rtype: List[str]
        """

        if len(digits) == 1:
            return self.m[digits]

        return self.multiply(self.m[digits[0]], self.letterCombinations(digits[1:]))


if __name__ == "__main__":
    print(Solution().letterCombinations("23"))
