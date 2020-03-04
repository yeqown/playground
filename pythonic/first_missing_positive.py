class Solution:

    @staticmethod
    def firstMissingPositive(nums:list):
        """
        :type nums: List[int]
        :rtype: int
        """
        n = len(nums)
        for i in range(n):
            while (nums[i] > 0 and nums[i] <= n and nums[nums[i]-1] != nums[i]):
                nums[nums[i]-1], nums[i] = nums[i], nums[nums[i]-1]

        print(nums)
        for i in range(n):
            if(nums[i] != i + 1):
                return i + 1
        return n + 1

if __name__ == "__main__":
    Solution.firstMissingPositive([7,8,9,11,12])
    Solution.firstMissingPositive([3,4,-1,1])
