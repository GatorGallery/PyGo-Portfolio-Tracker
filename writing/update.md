# Status Update by Noor Buchi and Anh Tran

## Description of the work already completed

As of writing this update, most of the project has
been completed. Both back-end and front-end
aspects are fully functional. Using Golang, we
used object oriented programming to implement
Portfolio and Stock objects with all needed
variables and operations. Functions such as buy,
sell, refresh, deposit, and withdraw can access,
modify, and update specific values in those
objects. Additionally, `StoreData()` writes the
contents of the object to a json file to be
retrieved later. All of the previously mentioned
operations are managed through `stockHandler.go`
which uses a command line interface to send
arguments. These functionalities are fully working
without any known errors. Additionally, the
front-end application and web interface are also
mostly completed. Implemented using Python, the
front-end uses Streamlit to create a user friendly
graphical interface for our application. The
Python modules use subprocess to build and run the
Go application. The previously mentioned command
line interface is also used in this process. In
addition to all of the basic operations done
through streamlit, we also use pandas, altair, and
matplotlib to create and visualize data frames.

## Description of the steps that you will take to finish the project

While the previous features are working, some of
them require further testing to ensure accuracy
and to get rid of any bugs resulting from unusual
input. This will be completed through both manual
testing of the web interface as well as potential
automated testing of the backend application. In
addition to establishing confidence in the
accuracy of our project, we plan to add other
features. One of the essential features on our
list involves storing all previous operations
previous and displaying then in the web interface
in a `History` section. This will require some
additions and modifications to our code but it
will be mostly straight forward. Other steps to
complete this project include making sure that our
code is well documented and easy to understand.
