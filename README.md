# GitHub Portfolio Generator with GPT-3.5

## Introduction

This GitHub Portfolio Generator is a tool that leverages the power of the GPT-3.5 language model to dynamically generate an HTML+CSS portfolio website based on your input. The generated portfolio showcases your github repos based on the list of names you provide.

## Features

- **Dynamic Content Generation**: Utilizes the GPT-3.5 model to dynamically generate content for your portfolio based on the information you provide.
- **Responsive Design**: The generated HTML+CSS code ensures a responsive and visually appealing design that looks great on various devices.
- **Easy Customization**: Easily customize the input parameters to tailor the portfolio to your preferences, such as projects, style.
- **Deeper Customization**: Consider adding additional data to the prompt, like a resume or social media links.
- **Deploy to Netlify**: Seamlessly deploy your generated portfolio to Netlify directly through their API, making it instantly accessible to the world.

## Prerequisites

Before using the GitHub Portfolio Generator, make sure you have the following:

- OpenAI API Key
- Netlify Account and token

## Usage

1. Clone this repository:

    ```bash
    git clone https://github.com/your-username/github-portfolio-generator.git
    ```

2. Make sure it works:

    ```bash
    cd github-portfolio-generator
    go build .
    ```

3. Set up environment variables:

    Create a `.env` file in the project root and add the following:

    ```env
    OPENAI_API_KEY=your_gpt3_api_key
    NETLIFY_ACCESS_TOKEN=your_netlify_api_key
    ```

5. Run generation + deployment:

    ```bash
    go run ./...
    ```

## Contributing

Feel free to contribute to this project by opening issues or submitting pull requests. Your feedback and improvements are highly appreciated!

## License

This GitHub Portfolio Generator is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.