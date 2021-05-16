package deploy

import "os"

/**
默认部署环境，使用命令行参数进行部署， --env=dev|test|prod, --set=xxx
*/
type DefaultDeploy struct {
}

func (d *DefaultDeploy) GetName() string {
	return "DefaultDeploy"
}

func (d *DefaultDeploy) Detect() *Info {
	properties := GetCommandLineProperties("")
	if envStr, ok := properties["env"]; ok {
		env := ParseEnv(envStr)
		set := properties["set"]
		return &Info{
			Env: env,
			Set: set,
		}
	}

	// 从环境变量中获取
	envStr := os.Getenv("env")
	if len(envStr) > 0 {
		env := ParseEnv(envStr)
		set := os.Getenv("set")
		return &Info{
			Env: env,
			Set: set,
		}
	}

	return &Info{
		Env: Dev,
	}
}
