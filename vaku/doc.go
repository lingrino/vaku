/*
Package vaku wraps an official Vault client with useful high-level functions.

Terminology

A 'path' in vault references a location. This location can either be a secret itself
or a folder containing secrets. A clean path does not start or end with a '/'.

A 'key' is the same as a path except that it can end in a '/', signifying for certain
that it is a folder. A trailing slash is a way of notifying that a path is a folder,
but lack of a trailing slash does not imply the opposite.

A 'folder' is a key that ends in a '/'. Folders are not secrets, but they contain
secrets and/or other folders.

Summary

Vaku is intended to provide useful functions for Vault that let users operate securely and efficiently. Using
the official Vault API and CLI it is only possible to take CRUD actions on individual paths, however normal
usage of vault will eventually lead to a desire for more advanced functions like the ability to copy or move
paths or even entire folders. The functions in this package can be broken into the distinct categroies below.

Helpers

These exported functions provide useful ways to manage Vault paths in code. They provide utilities
for cleaning, combining, and splitting Vault paths, keys, and folders. Many of these utilities exist
unexported in the official Vault Golang API, but they are useful enough to provide to end users here.

Path Functions

Path functions act on vault paths under both versions of the key/value secrets engine. They should not be
used on paths outside of those engines. The path functions in Vaku are opinionated and easy to use. However
this means that they do not support many of the features in the Vault API. For example, PathRead() does not
return any metadata or other information about a secret, only the data of the secret itself.

Folder Functions

Folder functions act on vault folders under both versions of the key/value secrets engine. They should
not be used on folders outside of those engines. In general, a folder function acts on all paths found
by listing the input path recursively. For example, FolderDelete() on "secret/test" will list all paths
within "secret/test" and its subfolders and then call PathDelete() on each one. Folder functions are
executed concurrently using a worker pool of goroutines and channels.
*/
package vaku
