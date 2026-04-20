<script lang="ts">
  import { page } from '$app/stores'
  import './layout.css'

  let { children } = $props()
  console.log($page.data)
   $inspect($page.data) 
</script>

{@render children()}
