# Questions for Shade


## 1. Generating random strings
- How do I generate a random password but not have to re-run/update pulumi components ? Similar to Terraform "random":

```
resource "random_string" "password2" {
  length  = 44
  special = false
}
```

Currently using Golang which re-creates above resources for every run:

```go
    tunnelSecret := utils.GenRandomString(44)

    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

    func randSeq(n int) string {
        b := make([]rune, n)
        for i := range b {
            b[i] = letters[rand.Intn(len(letters))]
        }
        return string(b)
    }

    func GenRandomString(n int) string {
        rand.Seed(time.Now().UnixNano())
        return randSeq(n)
    }
```
---

## 2. Leveraging TF provider for Boundary

- Here's the TF provider for boundary:

    | Terraform Provider| Link |
    | - | - |
    |Boundary | https://registry.terraform.io/providers/hashicorp/boundary/1.0.12|

- What's the best way to use from Pulumi ?

---

## 3. What's the best way to manage "variables" or mandatory inputs ? 

- How can I publish a "module" that require "inputs" or config ? 

- Like [Here](https://registry.terraform.io/modules/Bee-Projects/sfabric/aws/latest?tab=inputs) in a Terraform module I published a while ago.
- Are there stats on download count ?


![Inputs](./docs/TF_sfabric_inputs.png "blah")
