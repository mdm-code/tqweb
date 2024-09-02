/*
Package toml declares an adapter for external TOML packages to accommodate for
different implementations of data encoders and decoders. This limits the
necessary code changes when swapping out implementations. This can possibly
leave space for the user to configure the implementation they want to use.
*/
package toml
