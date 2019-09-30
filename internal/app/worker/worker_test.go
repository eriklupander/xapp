package worker

import (
	"bytes"
	"github.com/callistaenterprise/xapp/internal/app/filehandler/mock_filehandler"
	"github.com/callistaenterprise/xapp/internal/app/imageprocessor/mock_imageprocessor"
	"github.com/callistaenterprise/xapp/internal/app/persistence/mock_storage"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"testing"
)

func TestConsumeTweet(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Read test image from disk
	data, err := ioutil.ReadFile("../../../test/resources/testimage.png")
	assert.NoError(t, err)

	// Set up HTTP mocking
	gock.New("http://fake.url").
		Path("/path/123.png").
		Times(1).
		Reply(200).
		Body(bytes.NewBuffer(data))

	// Setup mocks and their expectations
	mockImgProc := mock_imageprocessor.NewMockImageProcessor(ctrl)
	mockImgProc.EXPECT().Hipsterize(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockDb := mock_storage.NewMockDatabase(ctrl)
	mockDb.EXPECT().ExistsByURL("http://fake.url/path/123.png").Return(false).Times(1)
	mockDb.EXPECT().Persist(gomock.Any()).Return(nil).Times(1)

	mockFh := mock_filehandler.NewMockFileHandler(ctrl)
	mockFh.EXPECT().Write(gomock.AssignableToTypeOf(""), gomock.Any()).Return(nil).Times(1)

	// Setup testee
	worker := NewTweetWorker(mockImgProc, mockDb, mockFh, nil)

	// Construct a tweet to test with
	tweet := &twitter.Tweet{
		Text: "My cat",
		User: &twitter.User{
			ScreenName: "Testsson",
		},
		Entities: &twitter.Entities{
			Media: []twitter.MediaEntity{
				{
					Type:     "photo",
					MediaURL: "http://fake.url/path/123.png",
				},
			},
		}}

	// When
	worker.processTweet(tweet)

	// Then - all expectations on mocks are asserted by gomock
	ctrl.Finish()
}
