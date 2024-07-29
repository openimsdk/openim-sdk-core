package vars

var (
	UserNum              int // user num
	SuperUserNum         int // number of users with all friends
	LargeGroupNum        int // number of large group
	CommonGroupNum       int // number of common group create by each user
	CommonGroupMemberNum int // number of common group member num
	SingleMessageNum     int // number of single message each user send
	GroupMessageNum      int // number of group message each user send

	ShouldRegister      bool // determine whether register
	ShouldImportFriends bool // determine whether import friends
	ShouldCreateGroup   bool // determine whether create group

	ShouldCheckGroupNum bool // determine whether check group num

	LoginRate float64 // number of login user rate
)
