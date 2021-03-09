# chrome-extension-downloader

Bulk download of Chrome-compatible Extensions from the web store.

## Usage

Just launch `chrome-ext-downloader.exe`
```
@noteworthy ➜  chrome-extension-downloader git(main) ✗  .\chrome-ext-downloader.exe

Enter list of extensions' URLs/IDs, then hit enter. Lines are separated with \n:

https://chrome.google.com/webstore/detail/hubspot-sales/oiiaigjnkhngdbnoookogelabohpglmd
https://chrome.google.com/webstore/detail/todoist-for-chrome/jldhpllghnbhlbpcmnajkpdmadaolakh
https://chrome.google.com/webstore/detail/simple-to-do-list/jimdhomgkpmmhhcegiebdajlkmjgikaf?hl=en

2021/03/09 16:13:29 Processing https://chrome.google.com/webstore/detail/hubspot-sales/oiiaigjnkhngdbnoookogelabohpglmd
2021/03/09 16:13:30 Processing https://chrome.google.com/webstore/detail/todoist-for-chrome/jldhpllghnbhlbpcmnajkpdmadaolakh
2021/03/09 16:13:30 Processing https://chrome.google.com/webstore/detail/simple-to-do-list/jimdhomgkpmmhhcegiebdajlkmjgikaf?hl=en
```

The extensions will be dropped in the same folder where the binary is located, both crx and zip.