Installation Instructions
     We strongly recommend installing Graphviz to its default location using the standard installer.

Default Path
    DEFAULT PATH: WINDOWS
        c:\Program Files\Graphviz*\dot.exe or c:\Program Files (x86)\Graphviz*\dot.exe
        depending on which version of Windows you are using.

Environment Variable
    If you have installed Graphviz somewhere other than the default location, you will need to define the environment variable GRAPHVIZ_DOT to point to the exact location of the DOT program. The variable must contain an executable, not a directory.
    ON WINDOWS: Create an environment variable GRAPHVIZ_DOT and point it to the DOT executable.
       GRAPHVIZ_DOT ==> For example: d:\example\dot.exe 