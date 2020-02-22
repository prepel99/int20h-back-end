This is backend side of project for mobile client: https://github.com/alibogzl/int20h_Quiz

We getting names of songs and their authors from audd.io API end then finding samples of these songs via deezer.com API

To ask for song you should make json POST request on  https://afternoon-gorge-54672.herokuapp.com/ with next body: {"text": "example of text"}

Then you get next model: 

```swift
struct Model: Codable {
        let songs: [Song]?
        let error: [String]?
    }
```
```swift
    struct Song: Codable {
        let title: String
        let preview: String
        let artist: Artist
    }
    
    struct Artist: Codable {
        let name: String
    }
```
