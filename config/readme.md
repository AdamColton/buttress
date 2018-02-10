## Gothic Config

A different approach to configs. One of the goals of gothic is to compile
everything into a single binary. That would include all configs so that a single
binary can run as any environment. Rather than keeping the configs in files,
the thinking is that the configs can be set in code.

### Recipes

#### Command line or os.Getenv
```Go
config.Environments("dev", "test", "prod")

if key := os.Getenv("key"){
	config.SetBytes("key").AsBase64(key)
} else {
	config.SetBytes("key").
		AsBase64("FEz_Oh8lFXY1u0ymBxipESwxlxprKUSFTHyPG5fGFt4=", "dev").
		AsBase64("JvQwfAAVxyn1WOt8NvQD4OJU--29mBKz8KoYGy1ghvs=", "test").
}
```

Where "key" is either an environment variable set for the user or is set from the
command line

```bash
key=ZoF4Z14lso3o4hlejIH_Q3TwayUxS5rmboiPLDH2Rs0= serviceName
```

This achieves two things. For security, the production key is not stored in the
code and must be set on the production machine. Second, we have default dev and
test keys, but we can also override them using an environment variable during
testing.

The same thing can be done with arguments passed in through the command line.