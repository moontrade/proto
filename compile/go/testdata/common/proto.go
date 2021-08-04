// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm

package common

type Direction byte

const (
	Direction_Down = Direction(0)

	Direction_Sideways = Direction(1)

	Direction_Up = Direction(2)
)

type Resolution int64

const (
	Resolution_S1 = Resolution(1000)

	Resolution_S3 = Resolution(3000)

	Resolution_S5 = Resolution(5000)

	Resolution_S10 = Resolution(10000)

	Resolution_S15 = Resolution(15000)

	Resolution_S30 = Resolution(30000)

	Resolution_M1 = Resolution(60000)

	Resolution_M2 = Resolution(120000)

	Resolution_M3 = Resolution(180000)

	Resolution_M4 = Resolution(240000)

	Resolution_M5 = Resolution(300000)

	Resolution_M6 = Resolution(360000)

	Resolution_M7 = Resolution(420000)

	Resolution_M8 = Resolution(480000)

	Resolution_M9 = Resolution(540000)

	Resolution_M10 = Resolution(600000)

	Resolution_M15 = Resolution(900000)

	Resolution_M20 = Resolution(1200000)

	Resolution_M30 = Resolution(1800000)

	Resolution_H1 = Resolution(3600000)

	Resolution_H2 = Resolution(7200000)

	Resolution_H3 = Resolution(10800000)

	Resolution_H4 = Resolution(14400000)

	Resolution_H6 = Resolution(21600000)

	Resolution_H8 = Resolution(28800000)

	Resolution_H10 = Resolution(36000000)

	Resolution_H12 = Resolution(43200000)

	Resolution_D = Resolution(86400000)

	Resolution_W = Resolution(604800000)
)
