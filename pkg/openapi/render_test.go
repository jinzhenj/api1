package openapi

import (
	"fmt"
	"testing"

	"github.com/api1/pkg/api1"
	"github.com/api1/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	t1 := `
group user

# Role
enum Role {

	# This is admin
	ADMIN

	# This is normal user
	USER
}

struct User {
	name: string
	age: int
	role: Role
}

interface user {

	# @route get /users
	listUsers(
		# Page Size from 0 to 100
		pageSize: int?,
		
		pageNo: int?,
		
		role: Role?
	): [User]

	# @route GET /users/:id
	getUser(id: int): User
}


# package some package

struct UserLoginRequest {
  # @go.validator: email
  email: string
  password: string
}

struct UserLoginResponse {
  token: string
}

struct UserRegisterRequest {
  email: string
  name: string
  password: string
}

struct ResetPasswordRequest {
  resetPasswordCode: string
  password: string
}

struct User2 {
  id: int
  name: string
  email: string
  role: string
}

interface UserController {
 
  # @route post /users/login
  userLogin(req: UserLoginRequest): UserLoginResponse

  # @route post /users/logout
  userLogout()

  # @route post /users/register
  userRegister(req: UserRegisterRequest)

  # @route get /users/activate
  userActivate(token: string)

  # @route get /users/reSendActivationEmail
  userReSendActivationEmail(email: string)

  # @route get /users/forgetPassword
  userForgetPassword(email: string)

  # @route post /users/reSetPassword
  userResetPassword(req: ResetPasswordRequest)

	# @route get /users
  getUserByCookie(): User?

  # set status = http.StatusConflict, if confliction detected.
	# @route get /users/check
  checkUserExists(email: string, name: string)
}

struct Like {
  like: boolean
  operated: boolean
}

interface DocController {

  # @route get /docs/:id
  getLikeCount(id: int): Like

  # @route put /docs/:id
  modifyLikeCount(id: int, like: boolean)
}

struct UploadFileResponse {
  key: string
  putUrl: string
  downloadUrl: string
  previewUrl: string
  headers: object
}

struct CreateFeedbackRequest {
  type: string     # 反馈类型
  module: [string] # 系统模块
  priority: string # 重要程度
  title: string    # 标题
  problemDataset: [int]  # 关联的数据集
  requestDataset: string # 扩充的数据集名称
  requestType: [string]  # 扩充类型
  content: string        # 反馈内容
  attachments: [string]  # 上传图片
}

interface FeedbackController {

  # @route post /file
  uploadFile(fileName: string, openRead: boolean?): UploadFileResponse

  # @route post /feedback
  createFeedback(req: CreateFeedbackRequest)
}

interface TypeController {

  # @route get /formatVersions
  listFormatVersions()
  
  # @route get /types
  listTypes()

  # @route put /types/:id
  # @go.middleware adminRequired
  modifyType(id: int)

  # @route post /types/batch
  # @go.middleware adminRequired
  batchCreateType()
  
  # @route delete /types/batch
  # @go.middleware adminRequired
  batchDeleteType()
}

interface DatasetController {

  # @route get /datasets
  listDatasets()

  # @route get /datasets/idName
  getDatasetIdAndName()

  # @route get /datasets/:id
  getDatasset(id: int)

  # @route get /datasets/:id/preview
  getDatasetPreview(id: int)

  # @route get /datasets/:id/files
  listDatasetFiles(id: int)

  # @route get /datasets/:id/sts
  getDatasetSts(id: int)

  # @route get /datasets/:id/similar
  getSimilarDatasets(id: int)

  # @route get /search/autoSuggest
  autoSuggest()

  # @route post /userEdit
  # @go.middleware loginRequired
  createUserEdit()

  # @route get /userEdit/:id
  # @go.middleware loginRequired
  getUserEdit(id: int)

  # @route get /datasets/operatorList
  searchDatasets()

  # @route get /datasets/:id/stats
  getDatasetStats(id: int)

  # @route post /datasets
  # @go.middleware adminRequired
  createDataset()
  
  # @route post /datasets/:id/actions/changeState
  # @go.middleware adminRequired
  changeDatasetState(id: int)

  # @route patch /datasets/:id
  # @go.middleware adminRequired
  modifyDataset(id: int)

  # @route delete /datasets/:id
  # @go.middleware adminRequired
  deleteDataset(id: int)

  # @route get /datasets/batch/preview
  # @go.middleware adminRequired
  batchImportDatasetPreview()
  
  # @route get /datasets/batch/import
  # @go.middleware adminRequired
  batchImportDataset()

  # @route get /datasets/:id/download/log
  logDownload(id: int)

  # @route get /datasets/:id/download
  getDatasetDownloadInfo(id: int)

  # @route put /datasets/:id/download
  modifyDatasetDownloadInfo(id: int)

  # @route get /dataset/download/log/export
  # @go.middleware adminRequired
  exportDownloadHistory()
}
 
interface CommentController {

#  	r.GET("/datasets/:id/comments", s.getComments)
# 	loginRequired := middleware.LoginRequired()
# 	r.POST("datasets/:id/comments", loginRequired, s.createComment)
# 	r.DELETE("/comments/:id", loginRequired, s.deleteComment)
}



interface sometest {
	# @route get /test
	test(matrix: [[int]])
}
`

	var err error

	if doc, err := parseAndRender(t1); err != nil {
		t.Fatalf("Parse error: %v", err)
	} else {
		fmt.Println(utils.ToJson(doc))
	}

	t10 := `
	  group t10
		
		interface T10 {
			# @route GET /test/{id}
			test(): int
		}
	`
	_, err = parseAndRender(t10)
	t.Log(err)
	assert.Error(t, err)

	t11 := `
	  group t11

		interface T11 {
			# @route POST /test
			test(param1: [int], param2: [[string]]): int
		}
	`
	_, err = parseAndRender(t11)
	t.Log(err)
	assert.Error(t, err)

	t12 := `
	  group t12

		interface T12 {
			# @route POST /test/:id
			test(id: int?): int
		}
	`
	_, err = parseAndRender(t12)
	t.Log(err)
	assert.Error(t, err)
}

func parseAndRender(s string) (*OpenAPI, error) {
	parser := api1.Parser{}
	schema, err := parser.Parse(s)
	if err != nil {
		return nil, err
	}

	render := Render{}
	return render.Render(schema)
}
