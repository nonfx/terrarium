**What is platform?**

A platform is a set of tools, services, and infrastructure that enables developers to build, deploy, and manage applications. It provides a foundation for software development by abstracting away the underlying infrastructure and providing a consistent set of APIs and services that developers can use to build their applications. Platforms can be internal or external, and can be hosted on-premises or in the cloud.

**What is platform engineering?**

Platform engineering is the discipline of designing and building toolchains, workflows, and platforms that enable self-service capabilities for software engineering organizations. It is an emerging field that focuses on enhancing developer productivity by reducing the complexity and uncertainty of modern software delivery. Platform engineering helps drive consistency and speed up common tasks by creating internal developer platforms (IDPs) that serve cross-company needs so that vertical application teams can serve their end-users. It is a process that organizations can use to leverage their cloud platform as efficiently as possible so that engineers can deliver value to production quickly and reliably.

Platform engineering is an evolution of business agility and DevOps, and it looks to break down more cross-organizational silos while reducing the cognitive load of application teams and duplicate work. It is usually achieved via the creation of the internal developer platform or IDP. Platform engineering is a powerful tool for scaling up software delivery processes without sacrificing quality, security, or efficiency. By simplifying and automating resource provisioning and management, it reduces operational complexity and removes friction from the development process.

In summary, platform engineering is the discipline of designing and building toolchains, workflows, and platforms that enable self-service capabilities for software engineering organizations. It helps drive consistency and speed up common tasks by creating internal developer platforms that serve cross-company needs so that vertical application teams can serve their end-users.

Further Reading:
[1] https://youtube.com/watch?v=Bfhl8kcSaEI
[2] https://platformengineering.org/blog/what-is-platform-engineering
[3] https://thenewstack.io/platform-engineering/
[4] https://circleci.com/blog/what-is-platform-engineering/
[5] https://www.liatrio.com/blog/what-is-platform-engineering-the-concept-behind-the-term
[6] https://humanitec.com/platform-engineering


**Why do we need platform engineering?**

Platform engineering is a discipline that focuses on enhancing developer productivity by reducing the complexity and uncertainty of modern software delivery. It is an emerging field that addresses some of the biggest challenges in software development, such as reducing operational complexity and removing friction from the development process. Platform engineering is a process that organizations can use to leverage their cloud platform as efficiently as possible so that engineers can deliver value to production quickly and reliably. It is a step forward in DevOps, enabling developers to follow DevOps practices more easily by creating a "golden path" that developers can use for rapid application development. Platform engineering is a powerful tool for scaling up software delivery processes without sacrificing quality, security, or efficiency. By simplifying and automating resource provisioning and management, it enables developers to ship more value to production faster. Some benefits of platform engineering include improving quality, delivering capabilities faster, and reducing cognitive load on development teams.


**Why IAC?**
Infrastructure as Code (IaC) is a method of managing and provisioning cloud infrastructure using programming techniques instead of manual processes. Here are some reasons why we need IaC:

1. **Automation**: IaC allows you to automate the process of provisioning and managing infrastructure. This means that you can create, modify, and delete infrastructure resources programmatically, which saves time and reduces the risk of human error.

2. **Consistency**: IaC ensures that your infrastructure is consistent across environments. By using the same code to provision infrastructure resources, you can ensure that they are configured in the same way every time. This reduces the risk of configuration drift and makes it easier to troubleshoot issues.

3. **Efficiency**: IaC makes it easier to deploy complex cloud architectures using a combination of pre-configured components. This reduces the amount of time and effort required to provision infrastructure resources, which increases efficiency and productivity for developers, architects, and administrators.

4. **Scalability**: IaC makes it easier to scale your infrastructure up or down as needed. By using code to provision infrastructure resources, you can quickly and easily add or remove resources as demand changes. This makes it easier to handle spikes in traffic or changes in workload.

5. **Security**: IaC can help improve security by making it easier to manage access controls and enforce security policies. By using code to provision infrastructure resources, you can ensure that security policies are applied consistently across your infrastructure.

6. **Documentation**: IaC makes it easier to document your infrastructure. By using code to provision infrastructure resources, you can create documentation that is easy to read and understand. This makes it easier to onboard new team members and troubleshoot issues.

In summary, IaC is important because it allows you to automate, audit, secure, and continuously deliver your infrastructure. It helps overcome common state management issues and ensures that your infrastructure is consistent, efficient, scalable, and secure.



**Why Terraform?**
Here are some reasons why Terraform is preferred over other IaC tools:

1. **Vendor-neutral**: Unlike other IaC tools, Terraform is vendor-neutral. You can use it to manage infrastructure in any supported platform or tool, such as Microsoft Azure, Google Cloud, AWS, Linode, and Oracle Cloud.

2. **Multi-cloud support**: Terraform supports multi-cloud deployments, which means that you can use it to manage infrastructure across various cloud providers. This makes it easier to manage infrastructure in a hybrid or multi-cloud environment.

3. **Declarative approach**: Terraform uses a declarative approach to infrastructure management, which means that you define the desired state of your infrastructure and let Terraform handle the details of how to achieve that state. This makes it easier to manage complex infrastructure configurations.

4. **Powerful configuration language**: Terraform uses the HashiCorp Configuration Language (HCL), which is concise and human-readable. This makes it easier to write and maintain Terraform code, even for complex infrastructure configurations.

5. **Large community**: Terraform has a large and active community of contributors, which means that you can find help and support easily. The community also contributes to the development of Terraform, which means that it is constantly improving.

6. **Extensive provider ecosystem**: Terraform has an extensive provider ecosystem that allows it to manage resources that may not be directly supported by other IaC tools like Chef, Ansible, Puppet, or CloudFormation.

7. **Consistent state management**: Terraform uses a state file to keep track of the current state of your infrastructure. This ensures that your infrastructure is consistent across environments and reduces the risk of configuration drift.

In summary, Terraform is preferred over other IaC tools because of its vendor-neutral approach, multi-cloud support, declarative approach, powerful configuration language, large community, extensive provider ecosystem, and consistent state management.

Further Reading:
[1] https://blog.gruntwork.io/why-we-use-terraform-and-not-chef-puppet-ansible-saltstack-or-cloudformation-7989dad2865c
[2] https://www.reddit.com/r/Terraform/comments/149bkxi/terraform_with_other_iac_tools/
[3] https://k21academy.com/terraform-iac/why-terraform-not-chef-ansible-puppet-cloudformation/
[4] https://spectralops.io/blog/terraform-vs-ansible/
[5] https://selleo.com/blog/why-choose-terraform-over-chef-puppet-ansible-saltstack-and-cloudformation
[6] https://ibatulanand.medium.com/extensive-comparison-of-iac-tools-49118e962ef8

**Why Adopt Platform-Based Infrastructure as Code (IaC)?**
Writing Platform-based IaC code has several benefits, including:
1. **Efficient Extensibility and Enhancement**
- **Initial Setup**: Requires some effort and resources to set up.
- **Long-term Benefits**: After the setup, extending and enhancing the platform becomes significantly easier and more efficient.
2. **Modular Terraform on Steroids**
- **Instant Code Generation**: Platform-based IaC enables quick and easy generation of the required Terraform code.
- **Next-Level Modularity**: Takes the concept of modular Terraform code to a more advanced and efficient level.
3. **Incorporation of Best Practices and Security**
- **Baked-In Security**: Helps in embedding best practices and security measures right from the start.
- **Proactive Stance**: Ensures a proactive approach to securing infrastructure and reducing vulnerabilities.
4. **Reduction in Time to Deliver Infrastructure Code**
- **Focus on Critical Areas**: Frees up time for the team to concentrate on improving site performance, enhancing security, and researching new tools and technologies.
- **Swift Project Delivery**: Contributes to faster and more efficient project completion and delivery.
5. **Enhanced Focus on Key Metrics**
- **Site Performance**: Allows teams to concentrate on improving the performance of the site.
- **Security Enhancements**: Provides time to focus on strengthening the security of the infrastructure.
- **Research on New Tools/Tech**: Encourages exploration and adoption of new tools and technologies.
6. **Ease of Onboarding for New DevOps Engineers**
- **Faster Integration**: Makes the onboarding process for new DevOps engineers quicker and more straightforward.
- **Empowerment**: Enables new team members to contribute effectively, even without extensive knowledge of Terraform.
7. **Promotion of Knowledge Democratization**
- **Knowledge Sharing**: Ensures that knowledge is accessible and shared among all team members, regardless of their experience level.
- **Collaborative Environment**: Fosters a work environment where all team members can actively participate and contribute to infrastructure development.

## Conclusion

Adopting a platform-based approach to Infrastructure as Code with Terraform is a strategic decision that brings efficiency, security, and best practices to the forefront of your DevOps and SRE operations. It empowers teams, streamlines processes, and ensures that your infrastructure is robust and ready for the challenges ahead.

Further Readings:
[1] https://www.forbes.com/sites/forbestechcouncil/2021/02/04/infrastructure-as-code-platform-what-it-is-and-how-it-can-be-used-in-multicloud-environments/?sh=266184a13831

** I already have modules, can I use them?**
If you already have modules, you can directly use them. We have written a guide on how to use existing modules in the platform-based code. You can find it here.


**How do I get started?**
You can get started by reading our documentation and following the steps outlined in the guide. You can also check out our examples to see how we have implemented platform-based IaC in real-world scenarios.

**I have a question that is not answered here. What should I do?**
Please refer to [support.md](https://raw.githubusercontent.com/cldcvr/terrarium/main/SUPPORT.md) Also join our Discord Community [here](Discord Link). We will be happy to answer any questions you may have.

**How do I contribute to this project?**
You can contribute to this project by submitting a pull request or opening an issue on GitHub. Please refer to our Review our [contributing guidelines](https://github.com/cldcvr/cldcvr-repo-template/blob/main/CONTRIBUTING.md).
 for more information.

**What is the license for this project?**
This project is licensed under the Apache License, Version 2.0. Please refer to the LICENSE file for more information.

**What is the roadmap for this project?**
We are currently working on adding more examples and documentation to this project. We also plan to add more features and functionality in the future. If you have any suggestions or feedback, please let us know.

**Why should I use this project?**
This project provides a set of tools and guidelines for writing platform-based IaC code. It also provides examples of how to use these tools and guidelines in real-world scenarios. If you are looking for a way to improve your infrastructure as code practices, this project can help you do that.

**Who is behind this project?**
This project is maintained by Ollion. We are a team of engineers who are passionate about DevOps. We have been using Terraform for several years now and have learned a lot along the way. We want to share our knowledge with the community so that others can benefit from it as well.

**When will this project be ready for production use?**
Stay tuned for updates on this project. We are working hard to make this project production-ready as soon as possible.
