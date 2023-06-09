package setupbooter

var SetupA = []PeerInstance{
	{ShouldConnectTo: 0}, // id 1 - should connect to id 0 (none)
	{ShouldConnectTo: 1}, // id 2 - should connect to id 1
	{ShouldConnectTo: 2}, // id 3 - should connect to id 2
} // Graph: 1 <- 2 <- 3

var SetupB = []PeerInstance{
	{ShouldConnectTo: 2}, // id 1 - should connect to id 2
	{ShouldConnectTo: 0}, // id 2 - should connect to id 0 (none)
	{ShouldConnectTo: 2}, // id 3 - should connect to id 2
} // Graph: 1 -> 2 <- 3

var SetupC = []PeerInstance{
	{ShouldConnectTo: 14},
	{ShouldConnectTo: 17},
	{ShouldConnectTo: 15},
	{ShouldConnectTo: 14},
	{ShouldConnectTo: 13},
	{ShouldConnectTo: 13},
	{ShouldConnectTo: 1},
	{ShouldConnectTo: 1},
	{ShouldConnectTo: 10},
	{ShouldConnectTo: 8},
	{ShouldConnectTo: 10},
	{ShouldConnectTo: 5},
	{ShouldConnectTo: 1},
	{ShouldConnectTo: 0},
	{ShouldConnectTo: 4},
	{ShouldConnectTo: 18},
	{ShouldConnectTo: 4},
	{ShouldConnectTo: 5},
	{ShouldConnectTo: 2},
	{ShouldConnectTo: 7},
}

// Graph for SetupC below:
/*
 *                                    ┌────┐
 *             ┌────┐     ┌───┐       │ 15 │◄───┬───┐
 *             │ 17 ├─────► 4 │◄──────┴────┘    │ 3 │
 *   ┌────┐    └─▲──┘     └─┬─┘                 └───┘
 *   │ 19 │      │          │                               ┌───┐
 *   └──┬─┘      │          │                               │ 9 │
 *      │      ┌─┴─┐        │        ┌───┐                  └─┬─┘
 *      └──────► 2 │     ┌──▼─┐◄─────┤ 1 ◄──┬───┐             │
 *             └───┘     │ 14 │      └─▲─┘  │   ├───┐         │
 * ┌────┐                └────┘        │    │   │ 8 │      ┌──▼─┐
 * │ 18 ├──────────────┐               │    │   └───┘◄─────┤ 10 │
 * └──▲─┘              │               │    │              └──▲─┘
 *    │              ┌─▼─┐     ┌────┬──┘    │                 │
 *    │              │ 5 ├────►│ 13 │     ┌─┴─┐               │
 *    └───┬────┐     └─▲─┘     └──▲─┘     │ 7 │               │
 *        │ 16 │       │          │       └───┘      ┌────┐   │
 *        └────┘       │          │         ▲        │ 11 ├───┘
 *                     │         ┌┴──┐      │        └────┘
 *                  ┌──┴─┐       │ 6 │      │
 *                  │ 12 │       └───┘      │
 *                  └────┘                 ┌┴───┐
 *                                         │ 20 │
 *                                         └────┘
 */

var SetupD = []PeerInstance{
	{ShouldConnectTo: 7},
	{ShouldConnectTo: 1},
	{ShouldConnectTo: 10},
	{ShouldConnectTo: 1},
	{ShouldConnectTo: 6},
	{ShouldConnectTo: 1},
	{ShouldConnectTo: 0},
	{ShouldConnectTo: 7},
	{ShouldConnectTo: 8},
	{ShouldConnectTo: 2},
}

// Graph for SetupD below:
/*
 * ┌───┐     ┌───┐
 * │ 9 ├─────► 8 │
 * └───┘     └─┬─┘
 *             │                               ┌───┐
 *             │                               │ 3 │
 *             │        ┌───┐                  └─┬─┘
 *           ┌─▼─┐◄─────┤ 1 ◄──┬───┐             │
 *           │ 7 │      └─▲─┘  │   ├───┐         │
 *           └───┘        │    │   │ 2 │      ┌──▼─┐
 *                        │    │   └───┘◄─────┤ 10 │
 *                        │    │              └────┘
 *      ┌───┐     ┌───┬───┘    │
 *      │ 5 ├────►│ 6 │      ┌─┴─┐
 *      └───┘     └───┘      │ 4 │
 *                           └───┘
 */
