package smartplaylists

import (
	"log"

	"github.com/zmb3/spotify"

	"github.com/josephbateh/senior-project-server/authentication"
	db "github.com/josephbateh/senior-project-server/database"
)

// PlaylistFromOtherPlaylists creates a new playlist from the tracks in the provided playlist
func PlaylistFromOtherPlaylists(userID string, name string, playlistIDs ...string) {
	// Get user from the DB
	user, err := db.GetUser(userID)
	if err != nil {
		log.Fatal(err)
	}

	// Get client from user
	client := authentication.GetClient(user.UserToken)

	// Get users playlists
	playlistPage, err := client.GetPlaylistsForUser(user.UserID)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the playlist name provided has already been used
	// If it hasn't, create the playlist
	playlists := playlistPage.Playlists
	var smartPlaylistID spotify.ID
	smartPlaylistIDSet := false
	for _, playlist := range playlists {
		pName := playlist.Name
		if name == pName {
			smartPlaylistID = playlist.ID
			smartPlaylistIDSet = true
		}
	}
	if !smartPlaylistIDSet {
		smartPlaylist, err := client.CreatePlaylistForUser(user.UserID, name, false)
		if err != nil {
			log.Fatal(err)
		}
		smartPlaylistID = smartPlaylist.ID
	}

	// Get all track IDs in one slice
	var tracks []spotify.ID
	for _, id := range playlistIDs {
		pTracks, err := client.GetPlaylistTracks(user.UserID, spotify.ID(id))
		if err != nil {
			log.Fatal(err)
		}
		for _, object := range pTracks.Tracks {
			tracks = append(tracks, object.Track.ID)
		}
	}

	// Clear playlist and add tracks
	tracksCurrentlyInPlaylist, _ := client.GetPlaylistTracks(user.UserID, smartPlaylistID)
	var currentTrackIDs []spotify.ID
	for _, object := range tracksCurrentlyInPlaylist.Tracks {
		currentTrackIDs = append(currentTrackIDs, object.Track.ID)
	}

	client.RemoveTracksFromPlaylist(user.UserID, smartPlaylistID, currentTrackIDs...)
	client.AddTracksToPlaylist(user.UserID, smartPlaylistID, tracks...)
}

func firstPlaylist() {
	// Get user ID
	userID := "jbspotifytest01"

	// Get user from DB
	user, err := db.GetUser(userID)
	if err != nil {
		log.Fatal(err)
	}

	// Get client from user
	client := authentication.GetClient(user.UserToken)

	// Get playlists
	playlistsPage, err := client.GetPlaylistsForUser(user.UserID)
	if err != nil {
		log.Fatal(err)
	}
	playlists := playlistsPage.Playlists

	// Get songs from playlist 1
	playlistOneTracks, err := client.GetPlaylistTracks(user.UserID, playlists[2].ID)
	if err != nil {
		log.Fatal(err)
	}
	playlistOneTrackObjects := playlistOneTracks.Tracks

	// Get songs from playlist 2
	playlistTwoTracks, err := client.GetPlaylistTracks(user.UserID, playlists[3].ID)
	if err != nil {
		log.Fatal(err)
	}
	playlistTwoTrackObjects := playlistTwoTracks.Tracks

	var tracksToBeAdded []spotify.ID
	// Create array that assigns to tracksToBeAdded
	for _, object := range playlistOneTrackObjects {
		tracksToBeAdded = append(tracksToBeAdded, object.Track.ID)
	}
	for _, object := range playlistTwoTrackObjects {
		tracksToBeAdded = append(tracksToBeAdded, object.Track.ID)
	}

	// Put songs from 1 and 2 into playlist 3
	_, err = client.AddTracksToPlaylist(user.UserID, playlists[0].ID, tracksToBeAdded...)
	if err != nil {
		log.Fatal(err)
	}
}
