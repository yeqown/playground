
class Solution:
    MAX = 999999999

    def minPathSum(self, grid: list[list[int]]) -> int:
        row = len(grid)
        if row == 0:
            return 0
        col = len(grid[0])
        if col == 0:
            return 0

        dp = [[0 for i in range(col)] for j in range(row)]
        # 初始化左上角
        dp[0][0] = grid[0][0]

        def load(x:int, y:int) -> int:
            if x < 0 or x >= col:
                return self.MAX
            if y < 0 or y >= row:
                return self.MAX
            return dp[x][y]
        
        for x in range(row):
            for y in range(col):
                dp = min(load(x-1,y),load(x,y-1)) + grid[x][y]
        
        print(dp)
        return grid[x][y]

