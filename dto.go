package golang_bittorrent_client

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"github.com/jackpal/bencode-go"
	"net"
	"os"
	"strconv"
	"time"
)

// File available as part of the torrent
type File struct {
	Length int    `bencode:"length"`
	Md5sum string `bencode:"md5sum"`
	Path   string `bencode:"path"`
}

// Data about the download itself
type MetaInfoData struct {
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Private     int    `bencode:"private"`
	Length      int    `bencode:"length"`
	Md5sum      string `bencode:"md5sum"`
	Files       []File `bencode:"files"`
}

// .torrent file description. Mostly metadata about the torrent
type MetaInfo struct {
	Announce     string       `bencode:"announce"`
	AnnounceList [][]string   `bencode:"announce-list"`
	Info         MetaInfoData `bencode:"info"`
	Encoding     string       `bencode:"encoding"`
	CreationDate int          `bencode:"creation date"`
	CreatedBy    string       `bencode:"created by"`
}

type Torrent struct {
	Path string
	Data MetaInfo
	Hash []byte
}

// NewTorrent builds a Torrent struct from the given .torrent file path
func NewTorrent(path string) (Torrent, error) {
	log.Debugf("Opening %s", path)
	file, err := os.Open(path)
	if err != nil {
		return Torrent{}, errors.New("Failed to open torrent file: " + err.Error())
	}
	defer file.Close()

	log.Debug("Decoding torrent file")
	info := MetaInfo{}
	err = bencode.Unmarshal(file, &info)
	if err != nil {
		return Torrent{}, errors.New("Failed to decode torrent file: " + err.Error())
	}

	return Torrent{Path: path, Data: info}, nil
}

// Calculates the unique identifier of the torrent
// by taking the SHA-1 hash of the torrent's 'info'; section
func computeInfoHash(torrentPath string) ([]byte, error) {

	file, err := os.Open(torrentPath)
	if err != nil {
		return nil, errors.New("Failed to open torrent: "; + err.Error())
	}
	defer file.Close()

	data, err := bencode.Decode(file)
	if err != nil {
		return nil, errors.New("Failed to decode torrent file: "; + err.Error())
	}

	torrentDict, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("Torrent file is not a dictionary")
	}

	infoBuffer := bytes.Buffer{}
	err = bencode.Marshal(&infoBuffer, torrentDict["info"])
	if err != nil {
		return nil, errors.New("Failed to encode info dict: " + err.Error())
	}

	hash := sha1.New()
	hash.Write(infoBuffer.Bytes())
	return hash.Sum(nil), nil
}

func NewTorrent(path string) (Torrent, error) {
	log.Debugf("Opening %s", path)
	file, err := os.Open(path)
	if err != nil {
		return Torrent{}, errors.New("Failed to open torrent file: " + err.Error())
	}
	defer file.Close()

	log.Debug("Decoding torrent file")
	info := MetaInfo{}
	err = bencode.Unmarshal(file, &info)
	if err != nil {
		return Torrent{}, errors.New("Failed to decode torrent file: " + err.Error())
	}

	log.Debug("Computing torrent info hash")
	infoHash, err := computeInfoHash(path)
	if err != nil {
		return Torrent{}, errors.New("Failed to compute info hash: " + err.Error())
	}

	log.Debugf("Announce URL: %s", info.Announce)
	log.Debugf("Hash: %x", infoHash)

	return Torrent{Path: path, Data: info, Hash: infoHash}, nil
}

type ClientEvent string
const (
	Started   ClientEvent = "started"
	Stopped               = "stopped"
	Completed             = "completed"
)

type TrackerRequest struct {
	InfoHash   []byte
	PeerId     []byte
	Port       int
	Uploaded   int
	Downloaded int
	Left       int
	Compact    int
	Event      ClientEvent
}


func generatePeerId() []byte {
	hash := sha1.New()
	// Current time
	hash.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))

	// Process ID
	hash.Write([]byte(strconv.Itoa(os.Getpid())))

	return hash.Sum(nil)
}

func (c *BittorentNetwork) Start(callback func(net.Conn)) (int, error) {
	var port = -1
	for i := 6881; i < 6890; i++ {
		log.Debug("Attempting to listen on port %d", i)
		listen, err := net.Listen("tcp", ":"+strconv.Itoa(i))
		if err == nil {
			c.listener = listen
			port = i
			break
		}
	}

	if c.listener == nil {
		return 0, errors.New("Unable to bind to port between 6881 and 6889")
	}

	log.Infof("Listening on port %d", port)
	go c.listenOnPort(callback)

	return port, nil
}

// First check the announce list
if len(torrent.Data.AnnounceList) > 0 {
for _, tier := range torrent.Data.AnnounceList {
// Randomly select trackers from each tier
for _, index := range rand.Perm(len(tier)) {
tracker := Tracker{tier[index]}
log.Debugf("(%s) Announcing...", tracker.Url)
resp, err := tracker.Announce(trackerRequest)
if err == nil && len(resp.FailureReason) == 0 {
return resp, true
}
log.Debugf("(%s) Announce failed: .", tracker.Url)
}
}
}

// If no announce list tracker was successful,
// try using the top level announce
tracker := &Tracker{Url: torrent.Data.Announce}
log.Debugf("(%s) Announcing...", tracker.Url)
resp, err := tracker.Announce(trackerRequest)
if err == nil {
return resp, true
}

log.Debugf("(%s) Announce failed.", tracker.Url)
return TrackerResponse{}, false



urlObj, err := url.Parse(c.Url)

if err != nil {
return TrackerResponse{}, fmt.Errorf("Unable to parse url: %s"+err.Error(), c.Url)
}

values := url.Values{}
values.Add("info_hash", string(request.InfoHash))
values.Add("peer_id", string(request.PeerId))
values.Add("port", strconv.Itoa(request.Port))
values.Add("uploaded", strconv.Itoa(request.Uploaded))
values.Add("downloaded", strconv.Itoa(request.Downloaded))
values.Add("left", strconv.Itoa(request.Left))
values.Add("compact", strconv.Itoa(request.Compact))
values.Add("event", string(request.Event))

urlObj.RawQuery = values.Encode()

result, err := http.Get(urlObj.String())
if err != nil || result.StatusCode < 200 || result.StatusCode >= 300 {
return TrackerResponse{}, errors.New("Tracker announcement failed")
}

var response TrackerResponse
bencode.Unmarshal(result.Body, &response)

func (c *Tracker) ParsePeers(peers string) ([]Peer, error) {
	if len(peers)%6 != 0 {
		return nil, errors.New("Peer string length is not a multiple of 6")
	}

	peerCount := len(peers) / 6
	peerList := make([]Peer, peerCount)
	for i := 0; i < peerCount; i++ {
		ipBytes := peers[i*6 : (i*6)+6]
		ip := net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
		peer := Peer{Address: ip.String(), Port: int(binary.BigEndian.Uint16([]byte(ipBytes[4:6])))}

		peerList[i] = peer
	}

	return peerList, nil
}
