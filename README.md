# VitGo

**Vite + Go + No Deps**

> VitGo is a Go module that lets you serve your Vite project from a Go-based web server.

- Steps

1. Build your vite project

2. Specific where to find the `dist/` directory

3. Then the module figures out how to load the generated vite application into a web page.

## Installation

```bash
go get -u github.com/botwayorg/vitgo
```

## Getting It Into Your Go Project

The first requirement is to [use ViteJS's tooling](https://vitejs.dev/guide/#scaffolding-your-first-vite-project) for your JavaScript code. The easiest thing to do is either to start out this way, or to create a new project and move your files into the directory that Vite creates. Using NPM:

```bash
# npm
npm create vite@latest

# yarn
yarn create vite@latest

# pnpm
pnpm create vite@latest
```

You will need to position your source files and the generated `dist/` directory so Go can find your project, the `manifest.json` file that describes it, and the assets that Vite generates for you. You may need to change your `vite.config.js` file (`vite.config.ts` if you prefer using Typescript) to make sure the manifest file is generated as well. Here's what I'm using:

```typescript
/**
 * @type {import('vite').UserConfig}
 */
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
  build: {
    manifest: "manifest.json",
    rollupOptions: {
      input: {
        main: "src/main.ts",
      },
    },
  },
});
```

This, however, is more than you need. A minimal config file would be:

```typescript
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    manifest: "manifest.json",
  },
});
```

The essential piece here is the vue plugin (or whatever plugin you need instead for React, Preact or Svelte) and the `build.manifest` line, since `vitgo` needs the manifest file to be present in order to work correctly.

```go
package main

import (
    "embed"
    "html/template"
    "net/http"

    "github.com/botwayorg/vitgo"
)

//go:embed "dist"
var dist embed.FS

var vitGo *vitgo.VitGo

func main() {
    // Production configuration.
   config := &vitgo.ViteConfig{
	   Environment: "production",
	   AssetsPath:  "dist",
	   EntryPoint:  "src/main.js",
	   Platform:    "react",
	   FS:          os.DirFS("frontend"),
   }

    // Development configuration
   config := &vitgo.ViteConfig{
	   Environment: "development",
	   AssetsPath:  "frontend",
	   EntryPoint:  "src/main.js",
	   Platform:    "react",
	   FS:          os.DirFS("frontend"),
   }

    // Parse the manifest and get a struct that describes
    // where the assets are.
    vgo, err := vitgo.NewVitGo(config)
    if err != nil {
        // bail!
    }

    vitGo = vgo

    // and set up your routes and start your server....
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    // Now you can pass the vgo object to an HTML template
    ts, err := template.ParseFiles("path/to/your-template.tmpl")
    if err != nil {
  	    // better handle this...
    }

    ts.Execute(respWriter, vitGo)
}
```

You will also need to serve your javascript, css and images used by your javascript code to the web. You can use a solution like [`http.FileServer`](https://pkg.go.dev/net/http#FileServer), or the wrapper the library implements that configures this for you:

```go
// using the standard library's multiplexer:
mux := http.NewServeMux()

// Set up a file server for our assets.
fsHandler, err := vgo.FileServer()
if err != nil {
    log.Println("could not set up static file server", err)

    return
}

mux.Handle("/src/", fsHandler)
```

Some router implementations may alternatively require you to do something more like:

```go
// chi router
mux := chi.NewMux()

...

mux.Handle("/src/*", fsHandler)
```

YMMV :-)

## Templates

Your template gets the needed tags and links by declaring the vgo object in your template and calling RenderTags on, as so:

```HTML
<!doctype html>
<html lang="en">

{{ $vue := . }}
    <head>
        <meta charset="utf-8">
        <title>Home - Vue Loader Test</title>

        {{ if $vue }}
          {{ $vue.RenderTags }}
        {{ end }}
    </head>
    <body>
      <div id="app"></div>
    </body>
</html>
```

You should check that the vgo (`$vue` in our example) is actually defined as I do here, since it will be nil unless you inject it into your template.

## Configuration

VitGo is fairly smart about your Vite Javascript project, and will examine your package.json file on start up. If you do not override the standard settings in your vite.config.js file, `vitgo` will probably choose to do the appropriate thing.

As mentioned above, a ViteConfig object must be passed to the `NewVitGo()` routine, with anything you want to override. Here are the major fields and how to use them:

| Field               | Purpose                                                                                                                                       | Default Setting                                                                                                     |
| ------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| **Environment**     | What mode you want vite to run in.                                                                                                            | development                                                                                                         |
| **FS**              | A fs.Embed or fs.DirFS                                                                                                                        | none; required.                                                                                                     |
| **JSProjectPath**   | Path to your Javascript files                                                                                                                 | frontend                                                                                                            |
| **AssetPath**       | Location of the built distribution directory                                                                                                  | _Production:_ dist                                                                                                  |
| **Platform**        | Any platform supported by Vite. vue and react are known to work; other platforms _may_ work if you adjust the other configurations correctly. | Based upon your package.json settings.                                                                              |
| **EntryPoint**      | Entry point script for your Javascript                                                                                                        | Best guess based on package.json                                                                                    |
| **ViteVersion**     | Vite major version ("2" or "3")                                                                                                               | Best guess based on your package.json file in your project. If you want to make sure, specify the version you want. |
| **DevServerPort**   | Port the dev server will listen on; typically 3000 in version 2, 5173 in version 3                                                            | Best guess based on version                                                                                         |
| **DevServerDomain** | Domain serving assets.                                                                                                                        | localhost                                                                                                           |
| **HTTPS**           | Whether the dev server serves HTTPS                                                                                                           | false                                                                                                               |
