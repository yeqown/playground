package calculator_with_tpl

import (
	"strings"
	"testing"
)

func Test_tpl(t *testing.T) {
	userScaleAnswer := []IScorer{
		score{score: 10}, // question 1 得分
		score{score: 20}, // question 2 得分
		score{score: 30}, // question 3 得分
		score{score: 20}, // question 4 得分
		score{score: 10}, // question 5 得分
		score{score: 10}, // question 6 得分
	}

	// Q: 多维度数据如何计算？并应用动态评测结果？
	// Solution1: 内置多种计算规则，将需要应用的规则ID放到量表的扩展中？
	// Solution2: 由量表模块，硬编码或者配置
	// Solution3: 设计多维度的规则，并将计算结果的过程，使用golang模板的方式，然后注入维度得分获得结果，如下：
	scaleCalcScoreRules := []string{
		"dimension1_score = ((q1+q2+q6) * 0.3 + (q3+q4+q5*2) * 0.7) / 2",
		"dimension2_score = ((q1+q2+q6+q6) * 0.7 + (q3+q4+q5) * 0.3) / 2",
		// ... more
	}

	scores, err := calcScoreWithsRules(userScaleAnswer, scaleCalcScoreRules)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data := map[string]interface{}{
		"scores": scores, // scores = [27, 23]
	}

	result := applyScaleResult(tpl, data)
	// result = 得分：27， 23，属于多动症患者
	result = strings.Trim(result, "\n")
	t.Log(result)
}
