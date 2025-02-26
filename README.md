<h1 align="center">MumbleMates - P2P Chat</h1>
<p align="center">
    <img src="assets/MumbleMatesLogo.png" alt="Logo of MumbleMates" height="250px"/>
    <br>
    MumbleMates is a Peer-to-Peer Chat written in GO
</p>


<hr>

## Technical Documentation

The documentation is automatically created when changes are merged and are uploaded as artifacts.
The latest version is also available under the latest Github Action: [📄 Technical Documentation](https://github.com/TIATIP-24-A-a/MumbleMates/actions/workflows/technical-docs.yml?query=branch%3Amain+is%3Asuccess)

### Building yourself

1. Install tectonic (recommended)
    - Website: https://tectonic-typesetting.github.io/en-US/install.html
    - Scoop: 
        ```
        scoop install main/tectonic
        ```
    - Brew:
        ```
        brew install tectonic
        ```
2. Run the following command:
    ```
    tectonic ./docs/technical/main.tex
    ```

Instead of running the command, there is also the possibilty to use a Visual Studio Code Extension [LaTeX Workshop](https://marketplace.visualstudio.com/items?itemName=James-Yu.latex-workshop)

Since the documentation is written in TeX any LaTeX engine should work. Though [tectonic](tectonic-typesetting.github.io) is recommended since the environment is built on it and setup is fairly simple.
